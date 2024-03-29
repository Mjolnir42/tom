/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type tcpServer struct {
	pipe     chan IPFIXMessage
	mux      *ipfixMux
	listener net.Listener
	quit     chan interface{}
	exit     chan interface{}
	wg       sync.WaitGroup
	err      chan error
	pool     *sync.Pool
	lm       *lhm.LogHandleMap
	shutdown bool
}

func newTCPServer(conf config.IPDaemon, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*tcpServer, error) {
	s := &tcpServer{
		pipe: mux.pipe(`inTCP`),
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
		mux:  mux,
		pool: pool,
		lm:   lm,
	}

	var err error
	if s.listener, err = net.Listen(ProtoTCP, conf.ListenADDR); err != nil {
		return nil, fmt.Errorf("tcpServer/Listen: %w", err)
	}

	s.wg.Add(1)
	go s.serve()
	return s, nil
}

func (s *tcpServer) Err() chan error {
	return s.err
}

func (s *tcpServer) Exit() chan interface{} {
	return s.exit
}

func (s *tcpServer) serve() {
	defer s.wg.Done()
	log.Println(`tcpServer: start serving clients`)

	connections := make(chan net.Conn)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

	acceptloop:
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					break acceptloop
				default:
					s.err <- fmt.Errorf("tcpServer/Accept/fatal: %w", err)
					close(s.exit)
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
			s.wg.Add(1)
			go func() {
				log.Printf("tcpServer: accepted connection from: %s\n",
					conn.RemoteAddr().String(),
				)
				s.handleConnection(conn)
				s.wg.Done()
			}()
		case <-s.exit:
			log.Println(`tcpServer: goroutine indicated fatal error`)
			break serveloop
		case <-s.quit:
			log.Println(`tcpServer: received shutdown signal`)
			break serveloop
		}
	}
}

func (s *tcpServer) Stop() chan error {
	s.shutdown = true
	go func(e chan error) {
		close(s.quit)
		s.listener.Close()
		s.wg.Wait()
		close(e)
	}(s.Err())
	return s.Err()
}

func (s *tcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	raddr := net.IP{}
	if tcp, err := net.ResolveTCPAddr(conn.RemoteAddr().Network(), conn.RemoteAddr().String()); err != nil {
		s.lm.GetLogger(`error`).Errorln(err)
		return
	} else {
		raddr = net.ParseIP(tcp.IP.String())
	}

ReadLoop:
	for {
		select {
		case <-s.quit:
			break ReadLoop
		default:
			conn.SetDeadline(time.Now().Add(750 * time.Millisecond))

			scanner := bufio.NewScanner(conn)
			scanner.Split(split)
			scanner.Buffer(make([]byte, IPFIXMaxSize+1, IPFIXMaxSize+1), IPFIXMaxSize)

			for scanner.Scan() {
				frame := IPFIXMessage{
					raddr: &raddr,
					body:  s.pool.Get().([]byte),
				}
				i := copy(frame.body, scanner.Bytes())
				frame.body = frame.body[:i]
				// send via UDP, but discard if buffered channel is full
				go func() {
					select {
					case s.pipe <- frame:
					default:
						s.pool.Put(frame.body)
					}
				}()

				// refresh deadline after a read and s.quit has not
				// been closed yet
				select {
				case <-s.quit:
					log.Printf("tcpServer: forcing close on connection from: %s\n",
						conn.RemoteAddr().String(),
					)
					break ReadLoop
				default:
					conn.SetDeadline(time.Now().Add(750 * time.Millisecond))
				}
			}

			if err := scanner.Err(); err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					// conn.Deadline triggered
					continue ReadLoop
				} else if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					// net package triggered timeout
					continue ReadLoop
				} else if err != io.EOF {
					s.err <- fmt.Errorf("tcpServer/Datastream/Split: %w", err)
				}
			}
			// scanner finished without error or timeout -> received EOF and
			// connection is closed
			break ReadLoop
		}
	}
	conn.Close()
}

type tcpClient struct {
	pipe      chan IPFIXMessage
	mux       *ipfixMux
	ping      chan ping
	quit      chan interface{}
	wg        sync.WaitGroup
	err       chan error
	conn      net.Conn
	connected bool
	pool      *sync.Pool
	lm        *lhm.LogHandleMap
	conf      config.SettingsIPFIX
	client    config.IPClient
}

func newTCPClient(conf config.SettingsIPFIX, cl config.IPClient, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*tcpClient, error) {
	c := &tcpClient{
		mux:       mux,
		ping:      make(chan ping),
		quit:      make(chan interface{}),
		err:       make(chan error),
		connected: false,
		conf:      conf,
		client:    cl,
		pool:      pool,
		lm:        lm,
	}
	c.pipe = c.mux.pipe(`outTCP`)

	c.wg.Add(1)
	go c.run()
	return c, nil
}

func (c *tcpClient) Stop() chan error {
	go func(e chan error) {
		close(c.quit)
		c.wg.Wait()
		close(e)
	}(c.err)
	return c.err
}

func (c *tcpClient) Err() chan error {
	return c.err
}

func (c *tcpClient) Input() chan IPFIXMessage {
	return c.pipe
}

func (c *tcpClient) run() {
	defer c.wg.Done()

	c.wg.Add(1)
	go c.Reconnect()

dataloop:
	for {
		frame := IPFIXMessage{
			body: c.pool.Get().([]byte),
		}
		select {
		case <-c.quit:
			log.Println(`TLSClient: shutdown signal received`)
			if c.conn != nil {
				// might be before first established connection
				c.conn.Close()
			}
			break dataloop
		case <-c.ping:
			continue dataloop
		case frame = <-c.pipe:
			if !c.connected {
				select {
				case c.pipe <- frame:
				default:
					// discard data while not connected and buffer is full
				}
				time.Sleep(125 * time.Millisecond)
				continue dataloop
			}

			if n, err := c.conn.Write(
				frame.body,
			); err != nil {
				c.err <- fmt.Errorf("TLSClient/Write: %w", err)
				c.connected = false
				c.conn.Close()
				select {
				case c.pipe <- frame:
				default:
					// discard data while if buffer is full
				}
			} else if n != len(frame.body) {
				c.connected = false
				c.conn.Close()
			}
			frame.body = frame.body[:0]
			c.pool.Put(frame.body)
		}
	}
}

func (c *tcpClient) Reconnect() {
	defer c.wg.Done()

	select {
	case <-c.quit:
		return
	default:
	}

	if c.conn != nil {
		c.conn.Close()
	}

connectloop:
	for ok := true; ok; ok = (c.connected == false) {
		dialer := &net.Dialer{
			Timeout:   750 * time.Millisecond,
			KeepAlive: 20 * time.Second,
		}
		var err error
		c.conn, err = dialer.Dial(ProtoTCP, c.client.ForwardADDR)
		if err != nil {
			c.err <- fmt.Errorf("TLSClient/Reconnect: %w", err)
			time.Sleep(time.Second)
			select {
			case <-c.quit:
				return
			default:
				continue connectloop
			}
		}
		log.Printf("TLSClient: connected to %s\n", c.client.ForwardADDR)
		break connectloop
	}

	c.connected = true
	c.ping <- ping{}

	c.wg.Add(1)
	go func() {
		readbuf := make([]byte, 512)

	detectloop:
		for {
			if err := c.conn.SetReadDeadline(
				time.Now().Add(250 * time.Millisecond),
			); err != nil {
				c.err <- fmt.Errorf("TLSClient/Reconnect: %w", err)
				c.connected = false
				c.conn.Close()
				break detectloop
			}
			if _, err := c.conn.Read(readbuf); err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					c.conn.SetReadDeadline(time.Time{})
				} else {
					c.err <- fmt.Errorf("TLSClient/Reconnect: %w", err)
					c.connected = false
					c.conn.Close()
					break detectloop
				}
			}
		}
		select {
		case <-c.quit:
			// intentional noop
		default:
			log.Printf("TLSClient: reconnecting to %s\n", c.client.ForwardADDR)
			c.wg.Add(1)
			c.Reconnect()
		}
		c.wg.Done()
	}()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
