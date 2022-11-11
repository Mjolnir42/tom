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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
