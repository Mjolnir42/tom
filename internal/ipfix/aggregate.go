/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"hash/fnv"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mjolnir42/flowdata"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

var twelveHours time.Duration = 12 * time.Hour

type procAggregate struct {
	conf config.SettingsIPFIX
	wg   sync.WaitGroup
	pool *sync.Pool
	lm   *lhm.LogHandleMap
	mux  *ipfixMux
	mtx  sync.RWMutex
	pipe chan flowdata.Record
	quit chan interface{}
	exit chan interface{}
	err  chan error
	agg  map[string]*Conn
}

func newAggregate(conf config.SettingsIPFIX, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*procAggregate, error) {
	m := &procAggregate{
		conf: conf,
		mux:  mux,
		pool: pool,
		lm:   lm,
		quit: make(chan interface{}),
		exit: make(chan interface{}),
		err:  make(chan error),
	}
	m.pipe = m.mux.flowdataPipe(`outAGG`)

	m.wg.Add(1)
	go m.run()
	return m, nil
}

func (m *procAggregate) run() {
	defer m.wg.Done()
	m.lm.GetLogger(`application`).Infoln(`ipfix.aggregation: starting input loop`)

runloop:
	for {
		select {
		case <-m.quit:
			m.lm.GetLogger(`application`).Infoln(`ipfix.aggregation: shutdown signal received`)
			break runloop
		case r := <-m.pipe:

			srcAddress := net.ParseIP(r.SrcAddress)
			dstAddress := net.ParseIP(r.DstAddress)

			c := &Conn{
				Protocol: r.Protocol,
			}

			// XXX TODO process record
			//
			// 2. check for local listening socket (https://github.com/cakturk/go-netstat)
			//    => if match, set remote port 0
			//    => if not match, set local port 0
			srcPort, dstPort := ephemeralPorts(r.SrcPort, r.DstPort)

			// srcAddress a, dstAddress b
			switch bytes.Compare([]byte(srcAddress), []byte(dstAddress)) {
			case 0:
				// a == b
				c.SrcAddress = srcAddress.String()
				c.SrcPort = srcPort
				c.DstAddress = dstAddress.String()
				c.DstPort = dstPort
			case -1:
				// a < b
				c.SrcAddress = srcAddress.String()
				c.SrcPort = srcPort
				c.DstAddress = dstAddress.String()
				c.DstPort = dstPort
			case +1:
				// a > b
				c.SrcAddress = dstAddress.String()
				c.SrcPort = dstPort
				c.DstAddress = srcAddress.String()
				c.DstPort = srcPort
			}

			if _, ok := m.agg[c.String()]; !ok {
				m.agg[c.String()] = c
			}
			m.agg[c.String()].Seen()
		}
	}
}

type Conn struct {
	SrcAddress string    `json:"srcAddr"`
	SrcPort    int       `json:"srcPort"`
	DstAddress string    `json:"dstAddr"`
	DstPort    int       `json:"dstPort"`
	Protocol   string    `json:"proto"`
	LastSeen   time.Time `json:"-"`
	CacheUntil time.Time `json:"-"`
	hash       string    `json:"-"`
}

func (c *Conn) String() string {
	if c.hash != `` {
		return c.hash
	}

	h := fnv.New128a()
	h.Write([]byte(c.SrcAddress))
	h.Write([]byte(strconv.Itoa(c.SrcPort)))
	h.Write([]byte(c.DstAddress))
	h.Write([]byte(strconv.Itoa(c.DstPort)))
	h.Write([]byte(c.Protocol))
	c.hash = hex.EncodeToString(h.Sum(nil))
	return c.hash
}

func (c *Conn) Seen() {
	c.LastSeen = time.Now().UTC()

	if c.CacheUntil.IsZero() {
		c.CacheUntil = c.LastSeen.Add(2 * twelveHours)
		return
	}

	if c.LastSeen.Add(twelveHours).After(c.LastSeen) {
		c.CacheUntil = c.CacheUntil.Add(twelveHours)
	}
}

func (c *Conn) Prune(now *time.Time) bool {
	if c.LastSeen.Add(twelveHours).Before(*now) {
		return true
	}
	return false
}

func (m *procAggregate) cron() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	select {
	case <-m.quit:
		return
	case <-m.exit:
		return
	case <-time.After(d):
	}

	tick := time.NewTicker(24 * time.Hour)
	<-m.prune()
	<-m.export()

	for {
		select {
		case <-m.quit:
			break
		case <-m.exit:
			break
		case <-tick.C:
			<-m.prune()
			<-m.export()
		}
	}
}

func (m *procAggregate) prune() chan interface{} {
	ret := make(chan interface{})
	go func(x chan interface{}) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		defer close(x)

		// XXX TODO
	}(ret)
	return ret
}

func (m *procAggregate) export() chan interface{} {
	ret := make(chan interface{})
	go func(x chan interface{}) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
		defer close(x)

		h, err := os.Hostname()
		if err != nil {
			m.err <- err
			return
		}
		h, _, _ = strings.Cut(h, `.`)

		req := proto.NewFlowRequest()
		req.Flow.Namespace = `fw-flow`
		req.Flow.Name = h
		req.Flow.Type = `slamdd-go/ipfix`

		for k, v := range m.agg {
			prop := proto.PropertyDetail{
				Attribute:  k + `_json`,
				ValidSince: v.LastSeen.Format(time.RFC3339),
				ValidUntil: v.CacheUntil.Format(time.RFC3339),
			}
			b, err := json.Marshal(v)
			if err != nil {
				m.err <- err
				continue
			}
			prop.Value = string(b)
			req.Flow.Property[prop.Attribute] = prop
		}

		spec := adm.Specification{
			Name: proto.CmdFlowEnsure,
			Body: req,
		}

		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet(``, flag.ContinueOnError), nil)

		err = adm.Perform(spec, ctx)
		if err != nil {
			m.err <- err
		}
	}(ret)
	return ret
}

func ephemeralPorts(src, dst uint16) (int, int) {
	switch {
	case src == 0:
		fallthrough
	case dst == 0:
		return int(src), int(dst)
	}

	var srcIsEph, dstIsEph bool
	if src >= 32768 {
		srcIsEph = true
	}
	if dst >= 32768 {
		dstIsEph = true
	}

	switch {
	case srcIsEph && !dstIsEph:
		return 0, int(dst)
	case !srcIsEph && dstIsEph:
		return int(src), 0
	default:
		return int(src), int(dst)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
