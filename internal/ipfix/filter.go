/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"sync"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type procFilter struct {
	conf    config.SettingsIPFIX
	quit    chan interface{}
	exit    chan interface{}
	err     chan error
	inpipe  chan []byte
	outpipe chan []byte
	mirror  chan []byte
	pool    *sync.Pool
	lm      *lhm.LogHandleMap
}

func newFilter(conf config.SettingsIPFIX, inpipe, outpipe, mirror chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*procFilter, error) {
	f := &procFilter{
		conf:    conf,
		quit:    make(chan interface{}),
		exit:    make(chan interface{}),
		err:     make(chan error),
		inpipe:  inpipe,
		outpipe: outpipe,
		mirror:  mirror,
		pool:    pool,
		lm:      lm,
	}
	// read from inpipe
	// filter
	// if processing contains aggregate
	//		copy to mirror
	// write to outpipe

	return f, nil
}

func (f *procFilter) serve() {
	defer f.wg.Done()
	f.lm.GetLogger(`application`).Infoln(`IPFIX|tlsServer: start serving clients`)

	connections := make(chan net.Conn)

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

	acceptloop:
		for {
			conn, err := f.listener.Accept()
			if err != nil {
				select {
				case <-f.quit:
					break acceptloop
				default:
					f.err <- fmt.Errorf("tlsServer/Accept/fatal: %w", err)
					close(f.exit)
					break acceptloop
				}
			}
			connections <- conn
		}
	}()

serveloop:
	for {
		select {
		case conn := <-connections:
			f.wg.Add(1)
			go func() {
				f.lm.GetLogger(`application`).Printf("IPFIX|tlsServer: accepted connection from: %s\n",
					conn.RemoteAddr().String(),
				)
				f.handleConnection(conn)
				f.wg.Done()
			}()
		case <-f.exit:
			f.lm.GetLogger(`application`).Println(`IPFIX|tlsServer: goroutine indicated fatal error`)
			break serveloop
		case <-f.quit:
			f.lm.GetLogger(`application`).Println(`IPFIX|tlsServer: received shutdown signal`)
			break serveloop
		}
	}
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
		f.listener.Close()
		f.wg.Wait()
		close(e)
	}(f.Err())
	return f.Err()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
