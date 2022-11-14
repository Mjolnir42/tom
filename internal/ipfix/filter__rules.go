/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mjolnir42/tom/internal/config"
)

func (f *procFilter) parseRules() error {
	f.parsedRules = make([]config.Rule, 0, len(f.conf.Filters.Rules))

	for i := range f.conf.Filters.Rules {
		f.conf.Filters.Rules[i] = strings.TrimSpace(
			f.conf.Filters.Rules[i],
		)
		tokens := strings.Split(f.conf.Filters.Rules[i], `;`)

		switch len(tokens) {
		case 3, 4:
		default:
			return fmt.Errorf(
				"Parsed rule has invalid length %d - (%s)",
				len(tokens),
				f.conf.Filters.Rules[i],
			)
		}

		r := config.Rule{}

		tokens[0] = strings.TrimSpace(tokens[0])
		switch tokens[0] {
		case `DROP`, `SET`, `REPLACE`, `PASS`:
			r.Action = tokens[0]
		default:
			return fmt.Errorf("unknown rule action: %s", tokens[0])
		}

		switch r.Action {
		case `DROP`, `SET`, `PASS`:
			if len(tokens) != 3 {
				return fmt.Errorf(
					"Parsed %s rule has invalid length %d - (%s)",
					tokens[0], len(tokens),
					f.conf.Filters.Rules[i],
				)
			}
		case `REPLACE`:
			if len(tokens) != 4 {
				return fmt.Errorf(
					"Parsed %s rule has invalid length %d - (%s)",
					tokens[0], len(tokens),
					f.conf.Filters.Rules[i],
				)
			}
		}

		tokens[1] = strings.TrimSpace(tokens[1])
		switch tokens[1] {
		case `ProtocolID`:
			r.MatchField = `ProtocolID`
			r.FieldType = `uint8`
		case `Protocol`:
			r.MatchField = `Protocol`
			r.FieldType = `string`
		case `TcpFlags`:
			r.MatchField = `TcpFlags`
			r.FieldType = `string`
		case `SrcAddress`:
			r.MatchField = `SrcAddress`
			r.FieldType = `string`
		case `SrcPort`:
			r.MatchField = `SrcPort`
			r.FieldType = `uint16`
		case `DstAddress`:
			r.MatchField = `DstAddress`
			r.FieldType = `string`
		case `DstPort`:
			r.MatchField = `DstPort`
			r.FieldType = `uint16`
		case `IPVersion`:
			r.MatchField = `IPVersion`
			r.FieldType = `uint8`
		case `AgentID`:
			r.MatchField = `AgentID`
			r.FieldType = `string`
		default:
			return fmt.Errorf("unknown match field: %s", tokens[1])
		}

		tokens[2] = strings.TrimSpace(tokens[2])
		if strings.HasPrefix(tokens[2], `NOT`) {
			r.InverseMatch = true
			tokens[2] = strings.TrimPrefix(tokens[2], `NOT`)
			tokens[2] = strings.TrimSpace(tokens[2])
		}
		matches := strings.Split(tokens[2], `,`)
		switch r.FieldType {
		case `uint8`:
			r.MatchValueUint8 = make([]uint8, len(matches))
			for i := range matches {
				num, err := strconv.ParseUint(matches[i], 10, 8)
				if err != nil {
					return err
				}
				r.MatchValueUint8[i] = uint8(num)
			}
		case `uint16`:
			r.MatchValueUint16 = make([]uint16, len(matches))
			for i := range matches {
				num, err := strconv.ParseUint(matches[i], 10, 16)
				if err != nil {
					return err
				}
				r.MatchValueUint16[i] = uint16(num)
			}
		case `string`:
			r.MatchValueString = matches
		}

		if r.Action == `REPLACE` {
			tokens[3] = strings.TrimSpace(tokens[3])
			switch r.FieldType {
			case `uint8`:
				num, err := strconv.ParseUint(tokens[3], 10, 8)
				if err != nil {
					return err
				}
				r.ReplaceValueUint8 = uint8(num)
			case `uint16`:
				num, err := strconv.ParseUint(tokens[3], 10, 16)
				if err != nil {
					return err
				}
				r.ReplaceValueUint16 = uint16(num)
			case `string`:
				r.ReplaceValueString = tokens[3]
			}

			switch r.MatchField {
			case `SrcAddress`, `SrcPort`, `DstAddress`, `DstPort`:
			default:
				return fmt.Errorf(
					"replace rule with unsupported match field: %s",
					r.MatchField,
				)
			}
		}
		if r.Action == `SET` {
			switch r.MatchField {
			case `AgentID`, `SrcAddress`, `DstAddress`:
			default:
				return fmt.Errorf(
					"set rule with unsupported match field: %s",
					r.MatchField,
				)
			}
		}
		if r.Action == `DROP` || r.Action == `PASS` {
			switch r.MatchField {
			case `AgentID`:
				return fmt.Errorf(
					"set rule with unsupported match field: %s",
					r.MatchField,
				)
			default:
			}
		}

		// store parsed rule
		f.parsedRules[i] = r
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
