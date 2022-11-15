/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"sync"

	flow "github.com/EdgeCast/vflow/ipfix"
	"github.com/mjolnir42/flowdata"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/vmware/go-ipfix/pkg/entities"
)

var (
	// templates memory cache
	mCache flow.MemCache
)

type procFilter struct {
	conf         config.SettingsIPFIX
	quit         chan interface{}
	exit         chan interface{}
	err          chan error
	wg           sync.WaitGroup
	pool         *sync.Pool
	lm           *lhm.LogHandleMap
	pipe         chan IPFIXMessage
	pipeConvert  chan IPFIXMessage
	pipeFilter   chan *MessagePack
	pipeOutput   chan *MessagePack
	outpipeIPFIX chan IPFIXMessage
	outpipeJSON  chan []byte
	outpipeFDR   chan flowdata.Record
	outpipeRawJS chan []byte
	fOutJSN      bool
	fRawJSN      bool
	fFmtJSN      string
	mux          *ipfixMux
	parsedRules  []config.Rule
	templates map[uint16]entities.Message
	sequences    map[uint32]uint32
}

func newFilter(conf config.SettingsIPFIX, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*procFilter, error) {
	f := &procFilter{
		conf: conf,
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
		mux:  mux,
		pool: pool,
		lm:   lm,
	}
	// read from inpipe
	// filter
	// if processing contains aggregate
	//		copy to mirror
	// write to outpipe
	for _, c := range conf.Clients {
		if !c.Enabled {
			continue
		}
		switch c.ForwardProto {
		case ProtoJSON:
			f.fOutJSN = true
			f.fRawJSN = c.Unfiltered
			f.fFmtJSN = c.Format
		}
	}

	if err := f.parseRules(); err != nil {
		return nil, err
	}

	// outFLT is the mux output to the filter
	f.pipe = f.mux.pipe(`outFLT`)

	f.pipeConvert = make(chan IPFIXMessage, 64)
	f.pipeFilter = make(chan *MessagePack, 64)
	f.pipeOutput = make(chan *MessagePack, 64)

	// inFLT is the mux input for filtered ipfix
	f.outpipeIPFIX = f.mux.pipe(`inFLT`)
	// inFLJ is the mux input for filtered JSON
	f.outpipeJSON = f.mux.jsonPipe(`inFLJ`)
	// inFLR is the mux input for filtered flowdata.Record
	f.outpipeFDR = f.mux.flowdataPipe(`inFLR`)
	// outJSN is the mux output for JSON data to the client
	f.outpipeRawJS = f.mux.jsonPipe(`outJSN`)

	// setup in-memory template cache
	mCache = flow.GetCache(f.conf.TemplFile)

	for i := 0; i < 8; i++ {
		f.wg.Add(1)
		go f.conversionWorker()

		f.wg.Add(1)
		go f.filterWorker()

		f.wg.Add(1)
		go f.outputWorker()
	}

	f.wg.Add(1)
	go f.run()
	return f, nil
}

func (f *procFilter) run() {
	defer f.wg.Done()
	f.lm.GetLogger(`application`).Infoln(`Filter module running`)

runloop:
	for {
		select {
		case <-f.quit:
			f.lm.GetLogger(`application`).Infoln(`Filter: shutdown signal received`)
			break runloop
		case frame := <-f.pipe:
			select {
			case f.pipeConvert <- frame:
			default:
			}
		}
	}
	// every 5 minutes => mCache.Dump(f.conf.TemplFile)
}

func (f *procFilter) Err() chan error {
	return f.err
}

func (f *procFilter) Exit() chan interface{} {
	return f.exit
}

func (f *procFilter) Stop() chan error {
	go func(e chan error) {
		close(f.quit)
		f.wg.Wait()
		close(e)
	}(f.Err())
	return f.Err()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
