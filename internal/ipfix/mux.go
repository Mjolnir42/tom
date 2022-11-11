/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"strings"
	"sync"

	"github.com/mjolnir42/flowdata"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type ipfixMux struct {
	conf config.SettingsIPFIX
	quit chan interface{}
	exit chan interface{}
	err  chan error
	wg   sync.WaitGroup
	// input channels from server listeners
	inUDP chan IPFIXMessage
	inTCP chan IPFIXMessage
	inTLS chan IPFIXMessage
	inJSN chan []byte
	// output channels to clients
	outUDP chan IPFIXMessage
	outTCP chan IPFIXMessage
	outTLS chan IPFIXMessage
	outJSN chan []byte
	// output and input channels to filter
	outFLT chan IPFIXMessage
	inFLT  chan IPFIXMessage
	inFLJ  chan []byte
	inFLR  chan flowdata.Record
	// output channel to TOM aggregator
	outAGG chan flowdata.Record
	// discard channels are served whenever an incorrect
	// channel is requested from a pipe() method
	discard     chan IPFIXMessage
	discardJSON chan []byte
	discardFDR  chan flowdata.Record

	pool        *sync.Pool
	lm          *lhm.LogHandleMap
	fOutUDP     bool
	fOutTCP     bool
	fOutTLS     bool
	fOutJSN     bool
	fRawUDP     bool
	fRawTCP     bool
	fRawTLS     bool
	fRawJSN     bool
	processing  bool
	filtering   bool
	aggregation bool
}

func newIPFIXMux(conf config.SettingsIPFIX, pool *sync.Pool, lm *lhm.LogHandleMap) *ipfixMux {
	m := &ipfixMux{
		conf: conf,
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
		pool: pool,
		lm:   lm,
		// input channels from server listeners
		inUDP: make(chan IPFIXMessage, 128), // input from udpServer
		inTCP: make(chan IPFIXMessage, 128), // input from tcpServer
		inTLS: make(chan IPFIXMessage, 128), // input from tlsServer
		inJSN: make(chan []byte, 128),       // input from jsonServer
		// output channels to clients
		outUDP: make(chan IPFIXMessage, 128), // output to udpClient
		outTCP: make(chan IPFIXMessage, 128), // output to tcpClient
		outTLS: make(chan IPFIXMessage, 128), // output to tlsClient
		outJSN: make(chan []byte, 128),       // output to jsonClient
		// output and input channels to filter
		outFLT: make(chan IPFIXMessage, 256),    // output to procFilter
		inFLT:  make(chan IPFIXMessage, 128),    // ipfix encoding
		inFLJ:  make(chan []byte, 128),          // JSON encoding
		inFLR:  make(chan flowdata.Record, 128), // flowdata encoding
		// output channel to TOM aggregator
		outAGG: make(chan flowdata.Record, 128), // flowdata encoding
		// discard channels
		discard:     make(chan IPFIXMessage, 2),
		discardJSON: make(chan []byte, 2),
		discardFDR:  make(chan flowdata.Record, 2),
	}

	for _, c := range conf.Clients {
		if !c.Enabled {
			continue
		}
		switch c.ForwardProto {
		case ProtoUDP:
			m.fOutUDP = true
			m.fRawUDP = c.Raw
		case ProtoTCP:
			m.fOutTCP = true
			m.fRawTCP = c.Raw
		case ProtoTLS:
			m.fOutTLS = true
			m.fRawTLS = c.Raw
		case ProtoJSON:
			m.fOutJSN = true
			m.fRawJSN = c.Raw
		}
	}
	if conf.Processing {
		m.processing = true
		for _, s := range strings.Split(conf.ProcessType, `,`) {
			switch s {
			case ProcFilter:
				m.filtering = true
			case ProcAggregate:
				m.aggregation = true
			}
		}
	}

	m.wg.Add(1)
	go m.run()
	return m
}

func (m *ipfixMux) run() {
	defer m.wg.Done()
	m.lm.GetLogger(`application`).Infoln(`mux: switching board running`)

	// if filtering is disabled, nobody is reading from outFLT and
	// writing back into inFLT.
	// the filtering module is also started if there is a JSON output
	if !m.filtering && !m.fOutJSN {
		m.wg.Add(1)
		go m.connectFilterChannel()
	}

	m.wg.Add(1)
	go m.connectJSONChannel()

runloop:
	for {
		select {
		case <-m.quit:
			m.lm.GetLogger(`application`).Infoln(`mux: shutdown signal received`)
			break runloop
		case frame := <-m.inUDP:
			select {
			case m.outFLT <- frame:
			default:
			}
		case frame := <-m.inTCP:
			select {
			case m.outFLT <- frame:
			default:
			}
		case frame := <-m.inTLS:
			select {
			case m.outFLT <- frame:
			default:
			}
		case buf := <-m.inJSN:
			select {
			case m.outJSN <- buf:
			default:
			}
		case frame := <-m.inFLT:
			if m.fOutUDP && !m.fRawUDP {
				go func(frame IPFIXMessage) {
					select {
					case m.outUDP <- frame:
					default:
					}
				}(frame.Copy())
			}
			if m.fOutTCP && !m.fRawTCP {
				go func(frame IPFIXMessage) {
					select {
					case m.outTCP <- frame:
					default:
					}
				}(frame.Copy())
			}
			if m.fOutTLS && !m.fRawTLS {
				go func(frame IPFIXMessage) {
					select {
					case m.outTLS <- frame:
					default:
					}
				}(frame.Copy())
			}
		case buf := <-m.inFLJ:
			if m.fOutJSN && !m.fRawJSN {
				go func(b []byte) {
					select {
					case m.outJSN <- b:
					default:
					}
				}(buf)
			}
		}
	}
}

func (m *ipfixMux) Err() chan error {
	return m.err
}

func (m *ipfixMux) Exit() chan interface{} {
	return m.exit
}

func (m *ipfixMux) Stop() chan error {
	go func(e chan error) {
		close(m.quit)
		m.wg.Wait()
		close(e)
	}(m.Err())
	return m.Err()
}

func (m *ipfixMux) pipe(p string) chan IPFIXMessage {
	switch p {
	case `inUDP`:
		return m.inUDP
	case `inTCP`:
		return m.inTCP
	case `inTLS`:
		return m.inTLS

	case `outUDP`:
		return m.outUDP
	case `outTCP`:
		return m.outTCP
	case `outTLS`:
		return m.outTLS

	case `outFLT`:
		return m.outFLT
	case `inFLT`:
		return m.inFLT

	default:
		return m.discard
	}
}

func (m *ipfixMux) jsonPipe(p string) chan []byte {
	switch p {
	case `inJSN`:
		return m.inJSN
	case `outJSN`:
		return m.outJSN
	case `inFLJ`:
		return m.inFLJ

	default:
		return m.discardJSON
	}
}

func (m *ipfixMux) flowdataPipe(p string) chan flowdata.Record {
	switch p {
	case `outAGG`:
		return m.outAGG
	default:
		return m.discardFDR
	}
}

func (m *ipfixMux) connectFilterChannel() {
	defer m.wg.Done()

	for {
		select {
		case <-m.quit:
			break
		case frame := <-m.outFLT:
			select {
			case m.inFLT <- frame:
			default:
			}
		}
	}
}

func (m *ipfixMux) connectJSONChannel() {
	defer m.wg.Done()

	for {
		select {
		case <-m.quit:
			break
		case buf := <-m.inJSN:
			select {
			case m.outJSN <- buf:
			default:
			}
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
