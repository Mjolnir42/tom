/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"sync"
)

type mockServer struct {
	quit     chan interface{}
	exit     chan interface{}
	err      chan error
	wg       sync.WaitGroup
	shutdown bool
}

func newMockServer() (*mockServer, error) {
	s := &mockServer{
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
	}

	s.wg.Add(1)
	go s.serve()
	return s, nil
}

func (s *mockServer) serve() {
	defer s.wg.Done()

DataLoop:
	for {
		select {
		case <-s.exit:
			break DataLoop
		case <-s.quit:
			break DataLoop
		}
	}
}

func (s *mockServer) Err() chan error {
	return s.err
}

func (s *mockServer) Exit() chan interface{} {
	return s.exit
}

func (s *mockServer) Stop() chan error {
	s.shutdown = true
	go func(e chan error) {
		close(s.quit)
		s.wg.Wait()
		close(e)
	}(s.err)
	return s.err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
