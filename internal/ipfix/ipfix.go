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
	udpServer     IPFIXServer
	udpClient     *udpClient
	tcpServer     IPFIXServer
	tcpClient     *tcpClient
	tlsServer     IPFIXServer
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

	pool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, IPFIXMaxSize, IPFIXMaxSize)
		},
	}

	mux := newIPFIXMux(conf.IPFIX, pool, lm)

	reg := registry{}

	// start client stage
	if conf.IPFIX.Forwarding {
		for i := range conf.IPFIX.Clients {
			if !conf.IPFIX.Clients[i].Enabled {
				continue
			}
			lm.GetLogger(`application`).Printf("IPFIX subsystem: starting forwarding client: %s", conf.IPFIX.Clients[i].ForwardProto)
			switch conf.IPFIX.Clients[i].ForwardProto {
			case ProtoUDP:
				reg.udpClient, err = newUDPClient(conf.IPFIX, conf.IPFIX.Clients[i], mux, pool, lm)
			case ProtoTCP:
				reg.tcpClient, err = newTCPClient(conf.IPFIX, conf.IPFIX.Clients[i], mux, pool, lm)
			case ProtoTLS:
				reg.tlsClient, err = newTLSClient(conf.IPFIX, conf.IPFIX.Clients[i], mux, pool, lm)
			case ProtoJSON:
			// TODO reg.jsonClient, err = newJSONClient(conf.IPFIX, outpipe, pool, lm)
			default:
				err = fmt.Errorf("Unsupported IPFIX output protocol: %s\n", conf.IPFIX.Clients[i].ForwardProto)
			}
			if err != nil {
				return
			}
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
				reg.procFilter, err = newFilter(conf.IPFIX, mux, pool, lm)
			case ProcAggregate:
				reg.procAggregate, err = newAggregate(conf.IPFIX, mux, pool, lm)
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
protoloop:
	for _, fixproto := range []string{ProtoUDP, ProtoTCP, ProtoTLS} {

		for _, srv := range conf.IPFIX.Servers {
			if fixproto == srv.ServerProto {
				if !srv.Enabled {
					continue
				}
				lm.GetLogger(`application`).Printf("IPFIX subsystem: starting server protocol %s", srv.ServerProto)
				switch srv.ServerProto {
				case ProtoUDP:
					reg.udpServer, err = newUDPServer(srv, mux, pool, lm)
				case ProtoTCP:
					reg.tcpServer, err = newTCPServer(srv, mux, pool, lm)
				case ProtoTLS:
					reg.tlsServer, err = newTLSServer(srv, mux, pool, lm)
				}
				if err != nil {
					return
				}
				continue protoloop
			}
		}
		lm.GetLogger(`application`).Printf("IPFIX subsystem: starting mock server protocol %s", fixproto)
		switch fixproto {
		case ProtoUDP:
			reg.udpServer, err = newMockServer()
		case ProtoTCP:
			reg.tcpServer, err = newMockServer()
		case ProtoTLS:
			reg.tlsServer, err = newMockServer()
		}
		if err != nil {
			return
		}
	}

	go func() {
		cancel := make(chan os.Signal, 1)
		signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)

		var ipfixSrvUDP, ipfixSrvTCP, ipfixSrvTLS IPFIXServer
		errordrain := make(chan error)

		for _, srv := range conf.IPFIX.Servers {
			switch srv.ServerProto {
			case ProtoUDP:
				ipfixSrvUDP = reg.udpServer
			case ProtoTCP:
				ipfixSrvTCP = reg.tcpServer
			case ProtoTLS:
				ipfixSrvTLS = reg.tlsServer
			}
		}
		drainUDP := make(chan interface{})
		drainTCP := make(chan interface{})
		drainTLS := make(chan interface{})
	runloop:
		for {
			select {
			case <-cancel:
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: received shutdown signal`)
				go chanCopy(ipfixSrvUDP.Stop(), errordrain, drainUDP)
				go chanCopy(ipfixSrvTCP.Stop(), errordrain, drainTCP)
				go chanCopy(ipfixSrvTLS.Stop(), errordrain, drainTLS)
				break runloop

			case <-ipfixSrvUDP.Exit():
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: UDP server process died`)
				lm.GetLogger(`error`).Infoln(`IPFIX subsystem: UDP server process died`)
				go chanCopy(ipfixSrvTCP.Stop(), errordrain, drainTCP)
				go chanCopy(ipfixSrvTLS.Stop(), errordrain, drainTLS)
				break runloop

			case <-ipfixSrvTCP.Exit():
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: TCP server process died`)
				lm.GetLogger(`error`).Infoln(`IPFIX subsystem: TCP server process died`)
				go chanCopy(ipfixSrvUDP.Stop(), errordrain, drainUDP)
				go chanCopy(ipfixSrvTLS.Stop(), errordrain, drainTLS)
				break runloop

			case <-ipfixSrvTLS.Exit():
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: TLS server process died`)
				lm.GetLogger(`error`).Infoln(`IPFIX subsystem: TLS server process died`)
				go chanCopy(ipfixSrvUDP.Stop(), errordrain, drainUDP)
				go chanCopy(ipfixSrvTCP.Stop(), errordrain, drainTCP)
				break runloop

			case <-mux.Exit():
				lm.GetLogger(`application`).Infoln(`IPFIX subsystem: mux switchboard died`)
				lm.GetLogger(`error`).Errorln(`IPFIX subsystem: mux switchboard died`)
				go chanCopy(ipfixSrvUDP.Stop(), errordrain, drainUDP)
				go chanCopy(ipfixSrvTCP.Stop(), errordrain, drainTCP)
				go chanCopy(ipfixSrvTLS.Stop(), errordrain, drainTLS)
				break runloop

			case err := <-ipfixSrvUDP.Err():
				lm.GetLogger(`error`).Errorln(err)

			case err := <-ipfixSrvTCP.Err():
				lm.GetLogger(`error`).Errorln(err)

			case err := <-ipfixSrvTLS.Err():
				lm.GetLogger(`error`).Errorln(err)
			}
		}

		lm.GetLogger(`application`).Infoln(`IPFIX subsystem: flushing pending errors during shutdown`)
		var servers uint8
		closed := false
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
			case <-drainUDP:
				servers = servers | 0b00000001
				if servers == 7 && !closed {
					close(errordrain)
					closed = true
				}
			case <-drainTCP:
				servers = servers | 0b00000010
				if servers == 7 && !closed {
					close(errordrain)
					closed = true
				}
			case <-drainTLS:
				servers = servers | 0b00000100
				if servers == 7 && !closed {
					close(errordrain)
					closed = true
				}
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

func chanCopy(i, o chan error, x chan interface{}) {
	for {
		select {
		case err := <-i:
			if err == nil {
				close(x)
				return
			}
			o <- err
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
