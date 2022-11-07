/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

const (
	IPFIXMaxSize         = 65535
	IPFIXVersion  uint16 = 10
	ProtoUDP             = `udp`
	ProtoTCP             = `tcp`
	ProtoTLS             = `tls`
	ProtoJSON            = `json`
	ProcFilter           = `filter`
	ProcAggregate        = `aggregate`
)

var (
	durationPlus5Min  time.Duration
	durationMinus5Min time.Duration
	ErrIncomplete     = errors.New("Incomplete IPFIX message")
)

func init() {
	var err error
	durationPlus5Min, err = time.ParseDuration(`+5m`)
	if err != nil {
		panic(err)
	}
	durationMinus5Min, err = time.ParseDuration(`-5m`)
	if err != nil {
		panic(err)
	}
}

type IPFIXServer interface {
	Err() chan error
	Exit() chan interface{}
	Stop() chan error
}

type registry struct {
	udpServer     *udpServer
	udpClient     *udpClient
	tcpServer     *tcpServer
	tcpClient     *tcpClient
	tlsServer     *tlsServer
	tlsClient     *tlsClient
	procFilter    *procFilter
	procAggregate *procAggregate
	//jsonClient *jsonClient
}

type ping struct{}

func New(conf config.SlamConfiguration, lm *lhm.LogHandleMap) (exit chan interface{}, err error) {
	exit = make(chan interface{})
	lm.GetLogger(`application`).Println(`IPFIX subsystem: starting`)
	if !conf.IPFIX.Enabled {
		lm.GetLogger(`application`).Println(`IPFIX subsystem: disabled by configuration`)
		close(exit)
		return
	}
	if !conf.IPFIX.Forwarding && !conf.IPFIX.Processing {
		lm.GetLogger(`application`).Errorln(`IPFIX receiving is activated, but both forwarding and processing disabled - skipping`)
		close(exit)
		return
	}

	// Message routing overview
	//
	// Enabled + Forwarding
	//    [server] -> outpipe -> [client]
	//
	// Enabled + Processing[Filter] + Forwarding
	//    [server] -> inpipe  -> [filter] -> outpipe -> [client]
	//
	// Enabled + Processing[Aggregate]
	//    [server] -> mirror -> [aggregate]
	//
	// Enabled + Processing[Filter,Aggregate] + Forwarding
	//    [server] -> inpipe  -> [filter]  -> outpipe -> [client]
	//                   \___> mirror  -> [aggregate]

	var pipe, inpipe, mirror, outpipe chan []byte
	if conf.IPFIX.Processing {
		inpipe = make(chan []byte, 1024)  // 1024 * 64 * 1024 = 64 MiB
		outpipe = make(chan []byte, 1024) // 1024 * 64 * 1024 = 64 MiB
		mirror = make(chan []byte, 1024)  // 1024 * 64 * 1024 = 64 MiB
		if conf.IPFIX.ProcessType == ProcAggregate {
			pipe = mirror
		} else {
			pipe = inpipe
		}
	} else {
		outpipe = make(chan []byte, 2048) // 2048 * 64 * 1024 = 128 MiB
		pipe = outpipe
	}

	pool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, IPFIXMaxSize, IPFIXMaxSize)
		},
	}
	reg := registry{}

	// start client stage
	if conf.IPFIX.Forwarding {
		lm.GetLogger(`application`).Println(`IPFIX subsystem: starting forwarding client`)
		switch conf.IPFIX.ForwardProto {
		case ProtoUDP:
			reg.udpClient, err = newUDPClient(conf.IPFIX, outpipe, pool, lm)
		case ProtoTCP:
			reg.tcpClient, err = newTCPClient(conf.IPFIX, outpipe, pool, lm)
		case ProtoTLS:
			reg.tlsClient, err = newTLSClient(conf.IPFIX, outpipe, pool, lm)
		case ProtoJSON:
		// TODO reg.jsonClient, err = newJSONClient(conf.IPFIX, outpipe, pool, lm)
		default:
			err = fmt.Errorf("Unsupported IPFIX output protocol: %s\n", conf.IPFIX.ServerProto)
		}
		if err != nil {
			return
		}
	} else {
		lm.GetLogger(`application`).Println(`IPFIX subsystem: forwarding disabled`)
	}

	// start processing stage
	if conf.IPFIX.Processing {
		for _, s := range strings.Split(conf.IPFIX.ProcessType, `,`) {
			lm.GetLogger(`application`).Println(`IPFIX subsystem: starting processing functions`)
			switch s {
			case ProcFilter:
				reg.procFilter, err = newFilter(conf.IPFIX, inpipe, outpipe, mirror, pool, lm)
			case ProcAggregate:
				reg.procAggregate, err = newAggregate(conf.IPFIX, mirror, pool, lm)
			default:
				err = fmt.Errorf("Unsupported IPFIX processing: %s\n", s)
			}
			if err != nil {
				return
			}
		}
	} else {
		lm.GetLogger(`application`).Println(`IPFIX subsystem: processing disabled`)
	}

	// start server stage
	lm.GetLogger(`application`).Printf("IPFIX subsystem: starting server protocol %s", conf.IPFIX.ServerProto)
	switch conf.IPFIX.ServerProto {
	case ProtoUDP:
		reg.udpServer, err = newUDPServer(conf.IPFIX, pipe, pool, lm)
	case ProtoTCP:
		reg.tcpServer, err = newTCPServer(conf.IPFIX, pipe, pool, lm)
	case ProtoTLS:
		reg.tlsServer, err = newTLSServer(conf.IPFIX, pipe, pool, lm)
	default:
		err = fmt.Errorf("Unsupported IPFIX input protocol: %s\n", conf.IPFIX.ServerProto)
	}
	if err != nil {
		return
	}

	go func() {
		cancel := make(chan os.Signal, 1)
		signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)

		var ipfixSrv IPFIXServer
		var errordrain chan error

		switch conf.IPFIX.ServerProto {
		case ProtoUDP:
			ipfixSrv = reg.udpServer
		case ProtoTCP:
			ipfixSrv = reg.tcpServer
		case ProtoTLS:
			ipfixSrv = reg.tlsServer
		}
	runloop:
		for {
			select {
			case <-cancel:
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: received shutdown signal`)
				errordrain = ipfixSrv.Stop()
				break runloop
			case <-ipfixSrv.Exit():
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: server process died`)
				lm.GetLogger(`error`).Infoln(`IPFIX subsystem: server process died`)
				break runloop
			case err := <-ipfixSrv.Err():
				lm.GetLogger(`error`).Errorln(err)
			}
		}

		lm.GetLogger(`application`).Infoln(`IPFIX subsystem: flushing pending errors during shutdown`)
	graceful:
		for {
			select {
			case <-time.After(time.Second * 15):
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: breaking graceful shutdown after 15s`)
				break graceful
			case err := <-errordrain:
				if err != nil {
					lm.GetLogger(`error`).Errorln(err)
					continue graceful
				}
				break graceful
			}
		}
		close(exit)
	}()

	return
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//  RFC 7011, 3.1 Message Header Format
	//
	//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |       Version Number          |            Length             |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |                           Export Time                         |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |                       Sequence Number                         |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |                    Observation Domain ID                      |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	// check there is at least the first half of the message header
	if len(data) < 8 {
		return 0, nil, nil
	}

	// check that octets 0 and 1 are the version number
	if binary.BigEndian.Uint16(data[:2]) != IPFIXVersion {
		// skip forward to seek the next valid header
		return 2, nil, nil
	}

	// read ipfix message length from octets 2, 3
	packLenByte := int(binary.BigEndian.Uint16(data[2:4]))

	// read export time from octets 4, 5, 6, 7
	unixT := int64(binary.BigEndian.Uint32(data[4:8]))
	// check that octets 4-7 can be interpreted as a timestamp of +/- 5
	// minutes from the local clock
	ts := time.Unix(unixT, 0).UTC()
	switch {
	case ts.Before(time.Now().UTC().Add(durationMinus5Min)):
		// skip forward to seek the next valid header
		fallthrough
	case ts.After(time.Now().UTC().Add(durationPlus5Min)):
		// skip forward to seek the next valid header
		return 8, nil, nil
	}

	// full ipfix data has not yet arrived
	if len(data) < packLenByte {
		switch atEOF {
		case false:
			return 0, nil, nil
		default:
			return 0, nil, ErrIncomplete
		}
	}

	return packLenByte, data[0:packLenByte], nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
