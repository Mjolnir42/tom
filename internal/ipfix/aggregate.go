/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"sync"

	"github.com/mjolnir42/flowdata"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type procAggregate struct {
	conf config.SettingsIPFIX
	wg   sync.WaitGroup
	pool *sync.Pool
	lm   *lhm.LogHandleMap
	mux  *ipfixMux
	pipe chan flowdata.Record
	quit chan interface{}
	exit chan interface{}
	err  chan error
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
		case <-m.pipe:
			// XXX TODO process record
			//
			// flowset-collection-name => hostname / type slamdd-go/ipfix
			// 1. normalize ordering src,dst
			// 2. check for local listening socket (https://github.com/cakturk/go-netstat)
			//    => if match, set remote port 0
			//    => if not match, set local port 0
			// 3. concat(srcaddress,dstaddress,srcport,dstport,proto) => asByte =>  hex-encode
			//    hexcode_json => attribute property
			// 4. set property validUntil +24h
			// 5. send proto.CmdFlowEnsure request to Tom
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
		}
	}
}

type Conn struct {
	SrcAddress string `json:"srcAddr"`
	SrcPort    int    `json:"srcPort"`
	DstAddress string `json:"dstAddr"`
	DstPort    int    `json:"dstPort"`
	Protocol   string `json:"proto"`
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
