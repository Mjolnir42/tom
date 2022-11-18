/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"

	"github.com/mjolnir42/flowdata"
)

type IPFIXMessage struct {
	raddr  *net.IP
	header IPFIXHeader
	body   []byte
}

type IPFIXHeader struct {
	Version    uint16
	Length     uint16
	ExportTime uint32
	SequenceNo uint32
	DomainID   uint32
	isTemplate bool
	numRecords uint32
	ClientID   uint32
}

type MessagePack struct {
	raddr   *net.IP
	header  IPFIXHeader
	ipfix   []byte
	records []*flowdata.Record
	jsons   [][]byte
}

func (i IPFIXMessage) Copy() IPFIXMessage {
	cc := IPFIXMessage{
		raddr:  i.raddr,
		header: i.header,
		body:   make([]byte, len(i.body)),
	}
	copy(cc.body, i.body)
	return cc
}

func (mp MessagePack) ExportIPFIX() IPFIXMessage {
	i := IPFIXMessage{
		raddr:  mp.raddr,
		header: mp.header,
		body:   make([]byte, len(mp.ipfix)),
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

func (mp MessagePack) SetClientID() {
	switch mp.raddr.To4() {
	case nil:
		// address is ip6
		r := []byte(mp.raddr.To16())
		b := make([]byte, 4, 4)
		copy(b[:2], r[2:4])
		copy(b[2:3], r[13:14])
		copy(b[3:], r[15:])
		mp.header.ClientID = binary.BigEndian.Uint32(b)
	default:
		// address is ip4
		s := strings.SplitN(mp.raddr.To4().String(), `.`, 4)
		b := make([]byte, 4, 4)
		for i := range b {
			num, _ := strconv.Atoi(s[i])
			b[i] = uint8(num)
		}
		mp.header.ClientID = binary.BigEndian.Uint32(b)
	}
	mp.header.ClientID = mp.header.ClientID | mp.header.DomainID
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
