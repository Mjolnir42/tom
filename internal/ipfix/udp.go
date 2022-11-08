/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type udpServer struct {
	listener   *net.UDPConn
	quit       chan interface{}
	exit       chan interface{}
	wg         sync.WaitGroup
	err        chan error
	remoteAddr string
	caFile     string
	pipe       chan []byte
	pool       *sync.Pool
	lm         *lhm.LogHandleMap
	shutdown   bool
}

func newUDPServer(conf config.IPDaemon, pipe chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*udpServer, error) {
	s := &udpServer{
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
		pipe: pipe,
		pool: pool,
		lm:   lm,
	}
	var err error
	var lUDPAddr *net.UDPAddr
	if lUDPAddr, err = net.ResolveUDPAddr(`udp`, conf.ListenADDR); err != nil {
		return nil, fmt.Errorf("udpServer/ResolveAddr: %w", err)
	}

	if s.listener, err = net.ListenUDP(`udp`, lUDPAddr); err != nil {
		return nil, fmt.Errorf("udpServer/ListenUDP: %w", err)
	}
	s.lm.GetLogger(`application`).Printf(
		"udpServer: listening on %s", conf.ListenADDR,
	)

	s.wg.Add(1)
	go s.serve()
	return s, nil
}

func (s *udpServer) serve() {
	defer s.wg.Done()

	buf := s.pool.Get().([]byte)
UDPDataLoop:
	for {
		select {
		case <-s.exit:
			s.lm.GetLogger(`error`).Println(`udpServer: goroutine indicated fatal error`)
			break UDPDataLoop
		case <-s.quit:
			s.lm.GetLogger(`application`).Println(`udpServer: received shutdown signal`)
			break UDPDataLoop
		default:
			s.listener.SetDeadline(time.Now().Add(750 * time.Millisecond))

			n, _, err := s.listener.ReadFromUDP(buf)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					// deadline triggered
					continue UDPDataLoop
				} else if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					// net package triggered timeout
					continue UDPDataLoop
				} else if errors.Is(err, net.ErrClosed) {
					// listener was closed
					if !s.shutdown {
						s.err <- fmt.Errorf("udpServer/ReadFromUDP/fatal: %w", err)
						close(s.exit)
					}
					break UDPDataLoop
				} else if err != io.EOF {
					s.err <- fmt.Errorf("udpServer/ReadFromUDP/fatal: %w", err)
					close(s.exit)
					break UDPDataLoop
				}
			}

			if n == 0 {
				// no data read, either with io.EOF or without
				continue UDPDataLoop
			}

			// make a data copy in an exact sized []byte
			data := make([]byte, n)
			copy(data, buf)

			select {
			case s.pipe <- data:
			default:
				// discard if buffered channel is full
			}
		}
	}
	//s.lm.GetLogger(`application`).Println(`UDP|Data: stopping client`)
	// XXX ch := client.Stop() -- XXX does this have to be replaced?
	/*
		drainloop:
			for {
				select {
				case e := <-ch:
					if e != nil {
						s.err <- e
						continue
					}
					// channel closed, read is nil
					break drainloop
				}
			}
	*/
	s.lm.GetLogger(`application`).Println(`UDP|Data: serve() done`)
}

func (s *udpServer) Err() chan error {
	return s.err
}

func (s *udpServer) Exit() chan interface{} {
	return s.exit
}

func (s *udpServer) Stop() chan error {
	s.shutdown = true
	go func(e chan error) {
		s.lm.GetLogger(`application`).Println(`UDP|STOP: Closing quit indicator channel`)
		close(s.quit)
		s.lm.GetLogger(`application`).Println(`UDP|STOP: closing listener`)
		s.listener.Close()
		s.lm.GetLogger(`application`).Println(`UDP|STOP: waiting for waitgroup`)
		s.wg.Wait()
		s.lm.GetLogger(`application`).Println(`UDP|STOP: closing error channel`)
		close(e)
		s.lm.GetLogger(`application`).Println(`UDP|STOP: done`)
	}(s.err)
	return s.err
}

type udpClient struct {
	inqueue chan []byte
	quit    chan interface{}
	err     chan error
	wg      sync.WaitGroup
	conf    config.SettingsIPFIX
	pool    *sync.Pool
	lm      *lhm.LogHandleMap
	UDPAddr *net.UDPAddr
	UDPConn *net.UDPConn
}

func newUDPClient(conf config.SettingsIPFIX, pipe chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*udpClient, error) {
	c := &udpClient{
		inqueue: pipe,
		quit:    make(chan interface{}),
		err:     make(chan error),
		conf:    conf,
		pool:    pool,
		lm:      lm,
	}

	var err error
	if c.UDPAddr, err = net.ResolveUDPAddr(ProtoUDP, c.conf.ForwardADDR); err != nil {
		return nil, fmt.Errorf("UDPClient/ResolveAddr: %w", err)
	}
	if c.UDPConn, err = net.DialUDP(ProtoUDP, nil, c.UDPAddr); err != nil {
		return nil, fmt.Errorf("UDPClient/Connect: %w", err)
	}

	c.wg.Add(1)
	go c.run()
	return c, nil
}

func (c *udpClient) Stop() chan error {
	go func(e chan error) {
		close(c.quit)
		c.wg.Wait()
		close(e)
	}(c.err)
	return c.err
}

func (c *udpClient) Err() chan error {
	return c.err
}

func (c *udpClient) Input() chan []byte {
	return c.inqueue
}

func (c *udpClient) run() {
	defer c.wg.Done()

runloop:
	for {
		msg := c.pool.Get().([]byte)
		select {
		case <-c.quit:
			c.lm.GetLogger(`application`).Println("UDPClient: received shutdown signal")
			break runloop

		case msg = <-c.inqueue:
		retryonerror:
			n, err := c.UDPConn.Write(msg)
			c.wg.Add(1)
			go func(e error) {
				defer c.wg.Done()
				if e != nil {
					c.err <- fmt.Errorf("UDPClient/Write: %w", e)
				}
			}(err)

		redial:
			if n != len(msg) {
				c.lm.GetLogger(`application`).Println("UDPClient: reconnecting after transient error....")
				// check if quit signal arrives during redial
				select {
				case <-c.quit:
					break runloop
				default:
				}
				time.Sleep(250 * time.Millisecond)

				// re-resolve UDP address
				if c.UDPAddr, err = net.ResolveUDPAddr(
					ProtoUDP, c.conf.ForwardADDR,
				); err != nil {
					c.wg.Add(1)
					go func(e error) {
						defer c.wg.Done()
						if err != nil {
							c.err <- fmt.Errorf("UDPClient/ResolveAddr: %w", e)
						}
					}(err)
					goto redial
				}
				// re-dial UDP connection
				if c.UDPConn, err = net.DialUDP(
					ProtoUDP, nil, c.UDPAddr,
				); err != nil {
					c.wg.Add(1)
					go func(e error) {
						defer c.wg.Done()
						if err != nil {
							c.err <- fmt.Errorf("UDPClient/Dial: %w", e)
						}
					}(err)
					goto redial
				}

				// retry sending current msg
				goto retryonerror
			}
			msg = msg[:0]
			c.pool.Put(msg)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
