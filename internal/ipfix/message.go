/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"net"
)

type IPFIXMessage struct {
	raddr *net.IP
	body  []byte
}

func (i IPFIXMessage) Copy() IPFIXMessage {
	cc := IPFIXMessage{
		raddr: i.raddr,
		body:  make([]byte, len(i.body)),
	}
	copy(cc.body, i.body)
	return cc
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
