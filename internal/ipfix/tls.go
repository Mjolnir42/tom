/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type tlsServer struct {
	listener    net.Listener
	certificate tls.Certificate
	quit        chan interface{}
	exit        chan interface{}
	wg          sync.WaitGroup
	err         chan error
	pipe        chan []byte
	pool        *sync.Pool
	lm          *lhm.LogHandleMap
}

func newTLSServer(conf config.IPDaemon, pipe chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*tlsServer, error) {
	s := &tlsServer{
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
		pipe: pipe,
		pool: pool,
		lm:   lm,
	}

	caPool := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(conf.CAFile); err != nil {
		return nil, fmt.Errorf("tlsServer/CA certificate: %w", err)
	} else {
		caPool.AppendCertsFromPEM(ca)
	}

	var err error
	if s.certificate, err = tls.LoadX509KeyPair(conf.CertFile, conf.CertKeyFile); err != nil {
		return nil, fmt.Errorf("tlsServer/CertificateKeypair: %w", err)
	}
	tlsConfig := &tls.Config{
		RootCAs:                  caPool,
		Time:                     time.Now().UTC,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP521},
		Certificates:             []tls.Certificate{s.certificate},
		ServerName:               conf.ServerName,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		ClientAuth: tls.NoClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	if s.listener, err = tls.Listen(`tcp`, conf.ListenADDR, tlsConfig); err != nil {
		return nil, fmt.Errorf("tlsServer/Listen: %w", err)
	}

	s.wg.Add(1)
	go s.serve()
	return s, nil
}

func (s *tlsServer) Err() chan error {
	return s.err
}

func (s *tlsServer) Exit() chan interface{} {
	return s.exit
}

func (s *tlsServer) serve() {
	defer s.wg.Done()
	s.lm.GetLogger(`application`).Infoln(`IPFIX|tlsServer: start serving clients`)

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
					s.err <- fmt.Errorf("tlsServer/Accept/fatal: %w", err)
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
				s.lm.GetLogger(`application`).Printf("IPFIX|tlsServer: accepted connection from: %s\n",
					conn.RemoteAddr().String(),
				)
				s.handleConnection(conn)
				s.wg.Done()
			}()
		case <-s.exit:
			s.lm.GetLogger(`application`).Println(`IPFIX|tlsServer: goroutine indicated fatal error`)
			break serveloop
		case <-s.quit:
			s.lm.GetLogger(`application`).Println(`IPFIX|tlsServer: received shutdown signal`)
			break serveloop
		}
	}
}

func (s *tlsServer) Stop() chan error {
	go func(e chan error) {
		close(s.quit)
		s.listener.Close()
		s.wg.Wait()
		close(e)
	}(s.Err())
	return s.Err()
}

func (s *tlsServer) handleConnection(conn net.Conn) {
	defer conn.Close()

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
				token := s.pool.Get().([]byte)
				i := copy(token, scanner.Bytes())
				// send via UDP, but discard if buffered channel is full
				go func() {
					select {
					case s.pipe <- token[:i]:
					default:
					}
					s.pool.Put(token)
				}()

				// refresh deadline after a read and s.quit has not
				// been closed yet
				select {
				case <-s.quit:
					s.lm.GetLogger(`application`).Printf("IPFIX|tlsServer: forcing close on connection from: %s\n",
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
					s.err <- fmt.Errorf("tlsServer/Datastream/Split: %w", err)
				}
			}
			// scanner finished without error or timeout -> received EOF and
			// connection is closed
			break ReadLoop
		}
	}
	conn.Close()
}

type tlsClient struct {
	inqueue   chan []byte
	ping      chan ping
	quit      chan interface{}
	wg        sync.WaitGroup
	err       chan error
	tlsConf   *tls.Config
	conn      *tls.Conn
	connected bool
	pool      *sync.Pool
	lm        *lhm.LogHandleMap
	conf      config.SettingsIPFIX
}

func newTLSClient(conf config.SettingsIPFIX, pipe chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*tlsClient, error) {
	c := &tlsClient{
		inqueue:   pipe,
		ping:      make(chan ping),
		quit:      make(chan interface{}),
		err:       make(chan error),
		connected: false,
		conf:      conf,
		pool:      pool,
		lm:        lm,
	}

	caPool := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(c.conf.CAFile); err != nil {
		return nil, fmt.Errorf("TLSClient/CA certificate: %w", err)
	} else {
		caPool.AppendCertsFromPEM(ca)
	}

	c.tlsConf = &tls.Config{
		RootCAs:          caPool,
		Time:             time.Now().UTC,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP521},
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
		},
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: false,
	}
	c.tlsConf.BuildNameToCertificate()

	c.wg.Add(1)
	go c.run()
	return c, nil
}

func (c *tlsClient) Stop() chan error {
	go func(e chan error) {
		close(c.quit)
		c.wg.Wait()
		close(e)
	}(c.err)
	return c.err
}

func (c *tlsClient) Err() chan error {
	return c.err
}

func (c *tlsClient) Input() chan []byte {
	return c.inqueue
}

func (c *tlsClient) run() {
	defer c.wg.Done()

	c.wg.Add(1)
	go c.Reconnect()

dataloop:
	for {
		msg := c.pool.Get().([]byte)
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
		case msg = <-c.inqueue:
			if !c.connected {
				select {
				case c.inqueue <- msg:
				default:
					// discard data while not connected and buffer is full
				}
				time.Sleep(125 * time.Millisecond)
				continue dataloop
			}

			if n, err := c.conn.Write(
				msg,
			); err != nil {
				c.err <- fmt.Errorf("TLSClient/Write: %w", err)
				c.connected = false
				c.conn.Close()
				select {
				case c.inqueue <- msg:
				default:
					// discard data while if buffer is full
				}
			} else if n != len(msg) {
				c.connected = false
				c.conn.Close()
			}
			msg = msg[:0]
			c.pool.Put(msg)
		}
	}
}

func (c *tlsClient) Reconnect() {
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
		c.conn, err = tls.DialWithDialer(dialer, ProtoTCP, c.conf.ForwardADDR, c.tlsConf)
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
		log.Printf("TLSClient: connected to %s\n", c.conf.ForwardADDR)
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
			log.Printf("TLSClient: reconnecting to %s\n", c.conf.ForwardADDR)
			c.wg.Add(1)
			c.Reconnect()
		}
		c.wg.Done()
	}()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
