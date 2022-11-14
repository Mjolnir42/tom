/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"net"

	"github.com/mjolnir42/flowdata"
)

type IPFIXMessage struct {
	raddr *net.IP
	body  []byte
}

type MessagePack struct {
	raddr   *net.IP
	ipfix   []byte
	records []*flowdata.Record
	jsons   [][]byte
}

func (i IPFIXMessage) Copy() IPFIXMessage {
	cc := IPFIXMessage{
		raddr: i.raddr,
		body:  make([]byte, len(i.body)),
	}
	copy(cc.body, i.body)
	return cc
}

func (mp MessagePack) ExportIPFIX() IPFIXMessage {
	i := IPFIXMessage{
		raddr: mp.raddr,
		body:  make([]byte, len(mp.ipfix)),
	}
	copy(i.body, mp.ipfix)
	return i
}

func (mp MessagePack) ExportJSON(s string) chan []byte {
	o := make(chan []byte)
	go func(pipe chan []byte, format string) {
		switch format {
		case `vflow`:
			pipe <- mp.jsons[0]
		case `flowdata`:
			for i := range mp.jsons {
				pipe <- mp.jsons[i]
			}
		}
		close(pipe)
	}(o, s)
	return o
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix