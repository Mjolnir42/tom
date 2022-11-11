/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"bytes"
	"sync"

	flow "github.com/EdgeCast/vflow/ipfix"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

var (
	// templates memory cache
	mCache flow.MemCache
)

type procFilter struct {
	conf    config.SettingsIPFIX
	quit    chan interface{}
	exit    chan interface{}
	err     chan error
	wg      sync.WaitGroup
	pool    *sync.Pool
	lm      *lhm.LogHandleMap
	fRawUDP bool
	fRawTCP bool
	fRawTLS bool
	fRawJSN bool
	mux     *ipfixMux
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
		case ProtoUDP:
			f.fRawUDP = c.Unfiltered
		case ProtoTCP:
			f.fRawTCP = c.Unfiltered
		case ProtoTLS:
			f.fRawTLS = c.Unfiltered
		case ProtoJSON:
			f.fRawJSN = c.Unfiltered
		}
	}

	mCache = flow.GetCache(f.conf.TemplFile)

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
		case frame := <-f.mux.Pipe(`outFLT`):

			if f.fRawUDP {
				select {
				case f.mux.Pipe(`outUDP`) <- frame.Copy():
				default:
				}
			}
			if f.fRawTCP {
				select {
				case f.mux.Pipe(`outTCP`) <- frame.Copy():
				default:
				}
			}
			if f.fRawTLS {
				select {
				case f.mux.Pipe(`outTLS`) <- frame.Copy():
				default:
				}
			}

			go func() {
				if outframe := f.filter(frame); outframe != nil {
					select {
					case f.mux.Pipe(`inFLT`) <- *outframe:
					default:
					}
				}
			}()
		}
	}
	// every 5 minutes => mCache.Dump(f.conf.TemplFile)
}

func (f *procFilter) filter(frame IPFIXMessage) *IPFIXMessage {
	decoder := flow.NewDecoder(*frame.raddr, frame.body)
	// mCache = ipfix.MemCache

	decodedMsg, err := decoder.Decode(mCache)
	if err != nil {
		f.lm.GetLogger(`error`).Println(`filter:`, err)
		if decodedMsg == nil {
			return nil
		}
	}
	buf := new(bytes.Buffer)

	if len(decodedMsg.DataSets) > 0 {
		b, err := decodedMsg.JSONMarshal(buf)
		if err != nil {
			f.lm.GetLogger(`error`).Println(`filter:`, err)
			return nil
		}
		f.lm.GetLogger(`request`).Println(string(b))
	}
	return &frame
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
