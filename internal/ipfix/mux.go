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

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type ipfixMux struct {
	conf        config.SettingsIPFIX
	quit        chan interface{}
	exit        chan interface{}
	err         chan error
	wg          sync.WaitGroup
	inUDP       chan IPFIXMessage
	inTCP       chan IPFIXMessage
	inTLS       chan IPFIXMessage
	inJSN       chan IPFIXMessage
	outUDP      chan IPFIXMessage
	outTCP      chan IPFIXMessage
	outTLS      chan IPFIXMessage
	outJSN      chan IPFIXMessage
	outFLT      chan IPFIXMessage
	inFLT       chan IPFIXMessage
	inFLJ       chan IPFIXMessage
	outAGG      chan IPFIXMessage
	discard     chan IPFIXMessage
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
		conf:    conf,
		quit:    make(chan interface{}),
		exit:    make(chan interface{}),
		err:     make(chan error),
		pool:    pool,
		lm:      lm,
		inUDP:   make(chan IPFIXMessage, 128), // input from udpServer
		inTCP:   make(chan IPFIXMessage, 128), // input from tcpServer
		inTLS:   make(chan IPFIXMessage, 128), // input from tlsServer
		inJSN:   make(chan IPFIXMessage, 128), // input from jsonServer
		outUDP:  make(chan IPFIXMessage, 128), // output to udpClient
		outTCP:  make(chan IPFIXMessage, 128), // output to tcpClient
		outTLS:  make(chan IPFIXMessage, 128), // output to tlsClient
		outJSN:  make(chan IPFIXMessage, 128), // output to jsonClient
		outFLT:  make(chan IPFIXMessage, 256), // output to procFilter
		inFLT:   make(chan IPFIXMessage, 128), // input from procFilter, ipfix encoding
		inFLJ:   make(chan IPFIXMessage, 128), // input from procFilter, JSON encoding
		outAGG:  make(chan IPFIXMessage, 128), // output to procAggregate
		discard: make(chan IPFIXMessage, 2),
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
	// writing back into inFLT
	if !m.filtering {
		m.wg.Add(1)
		go func() {
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
		}()
	}

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
		case frame := <-m.inJSN:
			select {
			case m.outFLT <- frame:
			default:
			}
			//case frame := <-m.inFLT: TODO

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

func (m *ipfixMux) Pipe(p string) chan IPFIXMessage {
	switch p {
	case `inUDP`:
		return m.inUDP
	case `inTCP`:
		return m.inTCP
	case `inTLS`:
		return m.inTLS
	case `inJSN`:
		return m.inJSN

	case `outUDP`:
		return m.outUDP
	case `outTCP`:
		return m.outTCP
	case `outTLS`:
		return m.outTLS
	case `outJSN`:
		return m.outJSN

	case `outFLT`:
		return m.outFLT
	case `inFLT`:
		return m.inFLT
	case `inFLJ`:
		return m.inFLJ

	case `outAGG`:
		return m.outAGG
	default:
		return m.discard
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
