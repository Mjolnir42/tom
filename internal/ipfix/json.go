/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type jsonServer struct {
	mux         *ipfixMux
	pipe        chan []byte
	listener    net.Listener
	certificate tls.Certificate
	quit        chan interface{}
	exit        chan interface{}
	wg          sync.WaitGroup
	err         chan error
	pool        *sync.Pool
	lm          *lhm.LogHandleMap
	user        string
	pass        string
}

func newJSONServer(conf config.IPDaemon, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*jsonServer, error) {
	s := &jsonServer{
		pipe: mux.jsonPipe(`inJSN`),
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
		mux:  mux,
		pool: pool,
		lm:   lm,
		user: conf.BasicUser,
		pass: conf.BasicPass,
	}

	switch {
	case s.user == ``:
		fallthrough
	case s.pass == ``:
		return nil, fmt.Errorf("jsonServer will not start without basic auth credentials")
	}

	caPool := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(conf.CAFile); err != nil {
		return nil, fmt.Errorf("jsonServer/CA certificate: %w", err)
	} else {
		caPool.AppendCertsFromPEM(ca)
	}

	var err error
	if s.certificate, err = tls.LoadX509KeyPair(conf.CertFile, conf.CertKeyFile); err != nil {
		return nil, fmt.Errorf("jsonServer/CertificateKeypair: %w", err)
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
		return nil, fmt.Errorf("jsonServer/Listen: %w", err)
	}

	s.wg.Add(1)
	go s.serve()
	return s, nil
}

func (s *jsonServer) Err() chan error {
	return s.err
}

func (s *jsonServer) Exit() chan interface{} {
	return s.exit
}

func (s *jsonServer) serve() {
	defer s.wg.Done()
	s.lm.GetLogger(`application`).Infoln(`IPFIX|tlsServer: start serving clients`)

	http.HandleFunc(`/submit`, s.inputHandler)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err := http.Serve(s.listener, nil)
		s.err <- err
		close(s.exit)
	}()

serveloop:
	for {
		select {
		case <-s.exit:
			s.lm.GetLogger(`application`).Println(`IPFIX|jsonServer: goroutine indicated fatal error`)
			break serveloop
		case <-s.quit:
			s.lm.GetLogger(`application`).Println(`IPFIX|jsonServer: received shutdown signal`)
			break serveloop
		case err := <-s.err:
			s.lm.GetLogger(`error`).Errorln(err)
		}
	}
}

func (s *jsonServer) Stop() chan error {
	go func(e chan error) {
		close(s.quit)
		s.listener.Close()
		s.wg.Wait()
		close(e)
	}(s.Err())
	return s.Err()
}

func (s *jsonServer) inputHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut, http.MethodPost, http.MethodPatch:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		return
	}

	auth := r.Header.Get(`Authorization`)
	if strings.HasPrefix(auth, `Basic`) {
		payload, err := base64.StdEncoding.DecodeString(auth[len(`Basic `):])
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 {
				switch {
				case string(pair[0]) != s.user:
				case string(pair[1]) != s.pass:
				default:
					goto authorizationSuccess
				}

			}
			http.Error(w, http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	}

authorizationSuccess:
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return

	}
}

type jsonClient struct {
	pipe      chan []byte
	mux       *ipfixMux
	ping      chan ping
	quit      chan interface{}
	wg        sync.WaitGroup
	err       chan error
	tlsConf   *tls.Config
	client    *resty.Client
	connected bool
	hostURL   string
	pool      *sync.Pool
	lm        *lhm.LogHandleMap
	conf      config.SettingsIPFIX
	clconf    config.IPClient
}

func newJSONClient(conf config.SettingsIPFIX, cl config.IPClient, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*jsonClient, error) {
	c := &jsonClient{
		mux:       mux,
		ping:      make(chan ping),
		quit:      make(chan interface{}),
		err:       make(chan error),
		connected: false,
		conf:      conf,
		clconf:    cl,
		pool:      pool,
		lm:        lm,
	}
	c.pipe = c.mux.jsonPipe(`outJSN`)

	host, err := url.Parse(`https://` + c.clconf.ForwardADDR + c.clconf.Endpoint)
	if err != nil {
		return nil, err
	}
	c.hostURL = host.String()

	c.client = resty.New().
		SetDisableWarn(true).
		SetHostURL(host.String())

	caPool := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(c.clconf.CAFile); err != nil {
		return nil, fmt.Errorf("jsonCLient/CA certificate: %w", err)
	} else {
		caPool.AppendCertsFromPEM(ca)
	}

	session := tls.NewLRUClientSessionCache(64)

	c.tlsConf = &tls.Config{
		RootCAs:          caPool,
		Time:             time.Now().UTC,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP521},
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
		},
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: false,
		ServerName:         strings.SplitN(host.Host, `:`, 2)[0],
		ClientSessionCache: session,
	}
	c.tlsConf.BuildNameToCertificate()

	c.client.SetTLSClientConfig(c.tlsConf).
		SetRootCertificate(c.clconf.CAFile)

	c.wg.Add(1)
	go c.run()
	return c, nil
}

func (c *jsonClient) Stop() chan error {
	go func(e chan error) {
		close(c.quit)
		c.wg.Wait()
		close(e)
	}(c.err)
	return c.err
}

func (c *jsonClient) Err() chan error {
	return c.err
}

func (c *jsonClient) Input() chan []byte {
	return c.pipe
}

func (c *jsonClient) run() {
	defer c.wg.Done()

dataloop:
	for {
		select {
		case <-c.quit:
			break dataloop
		case <-c.ping:
			continue dataloop
		case msg := <-c.pipe:
			r := c.client.R().
				SetBasicAuth(c.clconf.BasicUser, c.clconf.BasicPass).
				SetBody(msg).
				SetContentLength(true)

			var resp *resty.Response
			var err error
			switch c.clconf.Method {
			case http.MethodPut:
				resp, err = r.Put(c.hostURL)
			case http.MethodPatch:
				resp, err = r.Patch(c.hostURL)
			default:
				resp, err = r.Post(c.hostURL)
			}
			if err != nil {
				c.err <- err
				continue
			}
			if resp.IsError() {
				c.err <- fmt.Errorf(resp.Status())
			}
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
