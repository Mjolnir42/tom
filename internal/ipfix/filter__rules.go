/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"fmt"
	"strings"
	// "github.com/mjolnir42/flowdata"
)

type rule struct {
	MatchField         string
	FieldType          string
	MatchValueString   []string
	MatchValueUint8    []uint8
	MatchValueUint16   []uint16
	Action             string
	ReplaceValueString string
	ReplaceValueUint8  uint8
	ReplaceValueUint16 uint16
}

func (f *procFilter) parseRules() error {
	f.rules = make([]rule, 0, len(f.conf.Filters.Rules))

	for i := range f.conf.Filters.Rules {
		tok := strings.Split(f.conf.Filters.Rules[i], `;`)
		r := rule{}
		switch tok[0] {
		case `DROP`, `SET`, `REPLACE`:
			r.Action = tok[0]
		default:
			return fmt.Errorf("unknown rule action: %s", tok[0])
		}
		switch tok[1] {
		case `ProtocolID`:
		case `Protocol`:
		case `TcpFlags`:
		case `SrcAddress`:
		case `SrcPort`:
		case `DstAddress`:
		case `DstPort`:
		case `IPVersion`:
		case `AgentID`:
		default:
			return fmt.Errorf("unknown rule field: %s", tok[1])
		}
	}
	return nil
}

/*
DROP, ProtocolID, $value | NOT $value                                               uint8
DROP, Protocol, $value,$value | NOT $value,$value                                   string
DROP, TcpFlags, $value,$value | NOT $value,$value    SYN,FIN | NOT SYN,FIN          string
DROP, SrcAddress                                                                    string
DROP, SrcPort                                                                       uint16
DROP, DstAddress                                                                    string
DROP, DstPort                                                                       uint16
DROP, IPVersion                                                                     uint8


SET, AgentID, $value | clientIP
SET,

REPLACE, SrcAddress, $old, $new
REPLACE, SrcPort, $old, $new
REPLACE, DstAddress, $old, $new
REPLACE, DstPort, $old, $new
REPLACE, AgentID, $old, $new
*/

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
