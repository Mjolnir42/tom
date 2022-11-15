/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"github.com/mjolnir42/flowdata"
)

func (f *procFilter) filterWorker() {
	defer f.wg.Done()

loop:
	for {
		select {
		case <-f.quit:
			break loop
		case mp := <-f.pipeFilter:
			f.filter(mp)
		}
	}
}

func (f *procFilter) filter(pack *MessagePack) {
	f.applyRules(pack)

	// XXX TODO create IPFIX message into pack.ipfix
	//f.createIPFIX(&pack)

	f.pipeOutput <- pack
}

func (f *procFilter) applyRules(pack *MessagePack) {
	keep := make([]*flowdata.Record, 0, len(pack.records))

	// record index
recordloop:
	for ridx := range pack.records {

		// instruction index
	rulesloop:
		for iidx := range f.parsedRules {

			// check SET rules
			if f.parsedRules[iidx].Action == `SET` {
				switch f.parsedRules[iidx].MatchField {
				case `SrcAddress`:
					switch f.parsedRules[iidx].ReplaceValueString {
					case `clientIP`:
						pack.records[ridx].SrcAddress = pack.raddr.String()
					default:
						pack.records[ridx].SrcAddress = f.parsedRules[iidx].ReplaceValueString
					}
				case `DstAddress`:
					switch f.parsedRules[iidx].ReplaceValueString {
					case `clientIP`:
						pack.records[ridx].DstAddress = pack.raddr.String()
					default:
						pack.records[ridx].DstAddress = f.parsedRules[iidx].ReplaceValueString
					}
				case `AgentID`:
					switch f.parsedRules[iidx].ReplaceValueString {
					case `clientIP`:
						pack.records[ridx].AgentID = pack.raddr.String()
					default:
						pack.records[ridx].AgentID = f.parsedRules[iidx].ReplaceValueString
					}
				}
				continue rulesloop
			}

			// check REPLACE rules
			if f.parsedRules[iidx].Action == `REPLACE` {
				switch f.parsedRules[iidx].MatchField {
				case `SrcAddress`:
					for i := range f.parsedRules[iidx].MatchValueString {
						if pack.records[ridx].SrcAddress == f.parsedRules[iidx].MatchValueString[i] {
							pack.records[ridx].SrcAddress = f.parsedRules[iidx].ReplaceValueString
							continue rulesloop
						}
					}
					if f.parsedRules[iidx].InverseMatch {
						pack.records[ridx].SrcAddress = f.parsedRules[iidx].ReplaceValueString
						continue rulesloop
					}
				case `DstAddress`:
					for i := range f.parsedRules[iidx].MatchValueString {
						if pack.records[ridx].DstAddress == f.parsedRules[iidx].MatchValueString[i] {
							pack.records[ridx].DstAddress = f.parsedRules[iidx].ReplaceValueString
							continue rulesloop
						}
					}
					if f.parsedRules[iidx].InverseMatch {
						pack.records[ridx].DstAddress = f.parsedRules[iidx].ReplaceValueString
						continue rulesloop
					}
				case `SrcPort`:
					for i := range f.parsedRules[iidx].MatchValueUint16 {
						if pack.records[ridx].SrcPort == f.parsedRules[iidx].MatchValueUint16[i] {
							pack.records[ridx].SrcPort = f.parsedRules[iidx].ReplaceValueUint16
							continue rulesloop
						}
					}
					if f.parsedRules[iidx].InverseMatch {
						pack.records[ridx].SrcPort = f.parsedRules[iidx].ReplaceValueUint16
						continue rulesloop
					}
				case `DstPort`:
					for i := range f.parsedRules[iidx].MatchValueUint16 {
						if pack.records[ridx].DstPort == f.parsedRules[iidx].MatchValueUint16[i] {
							pack.records[ridx].DstPort = f.parsedRules[iidx].ReplaceValueUint16
							continue rulesloop
						}
					}
					if f.parsedRules[iidx].InverseMatch {
						pack.records[ridx].DstPort = f.parsedRules[iidx].ReplaceValueUint16
						continue rulesloop
					}
				}
			}

			// check DROP and PASS rules
			if f.parsedRules[iidx].Action == `DROP` || f.parsedRules[iidx].Action == `PASS` {
				match := false
				switch f.parsedRules[iidx].MatchField {
				case `ProtocolID`:
					for i := range f.parsedRules[iidx].MatchValueUint8 {
						if pack.records[ridx].ProtocolID == f.parsedRules[iidx].MatchValueUint8[i] {
							match = true
							break
						}
					}
				case `IPVersion`:
					for i := range f.parsedRules[iidx].MatchValueUint8 {
						if pack.records[ridx].IPVersion == f.parsedRules[iidx].MatchValueUint8[i] {
							match = true
							break
						}
					}
				case `SrcPort`:
					for i := range f.parsedRules[iidx].MatchValueUint16 {
						if pack.records[ridx].SrcPort == f.parsedRules[iidx].MatchValueUint16[i] {
							match = true
							break
						}
					}
				case `DstPort`:
					for i := range f.parsedRules[iidx].MatchValueUint16 {
						if pack.records[ridx].DstPort == f.parsedRules[iidx].MatchValueUint16[i] {
							match = true
							break
						}
					}
				case `Protocol`:
					for i := range f.parsedRules[iidx].MatchValueString {
						if pack.records[ridx].Protocol == f.parsedRules[iidx].MatchValueString[i] {
							match = true
							break
						}
					}
				case `SrcAddress`:
					for i := range f.parsedRules[iidx].MatchValueString {
						if pack.records[ridx].SrcAddress == f.parsedRules[iidx].MatchValueString[i] {
							match = true
							break
						}
					}
				case `DstAddress`:
					for i := range f.parsedRules[iidx].MatchValueString {
						if pack.records[ridx].DstAddress == f.parsedRules[iidx].MatchValueString[i] {
							match = true
							break
						}
					}
				case `TcpFlags`:
					for i := range f.parsedRules[iidx].MatchValueString {
						switch f.parsedRules[iidx].MatchValueString[i] {
						case `NS`, `ns`:
							if pack.records[ridx].TcpFlags.NS == true {
								match = true
								break
							}
						case `CWR`, `cwr`:
							if pack.records[ridx].TcpFlags.CWR == true {
								match = true
								break
							}
						case `ECE`, `ece`:
							if pack.records[ridx].TcpFlags.ECE == true {
								match = true
								break
							}
						case `URG`, `urg`:
							if pack.records[ridx].TcpFlags.URG == true {
								match = true
								break
							}
						case `ACK`, `ack`:
							if pack.records[ridx].TcpFlags.ACK == true {
								match = true
								break
							}
						case `PSH`, `psh`:
							if pack.records[ridx].TcpFlags.PSH == true {
								match = true
								break
							}
						case `RST`, `rst`:
							if pack.records[ridx].TcpFlags.RST == true {
								match = true
								break
							}
						case `SYN`, `syn`:
							if pack.records[ridx].TcpFlags.SYN == true {
								match = true
								break
							}
						case `FIN`, `fin`:
							if pack.records[ridx].TcpFlags.FIN == true {
								match = true
								break
							}
						}
					}
				}
				if f.parsedRules[iidx].InverseMatch && !match {
					match = true
				}
				if match {
					switch f.parsedRules[iidx].Action {
					case `DROP`:
						continue recordloop
					case `PASS`:
						r := pack.records[ridx].Copy()
						keep = append(keep, &r)
						continue recordloop
					}
				}
			}
		}
		// record passed all rules without hitting a DROP rule
		r := pack.records[ridx].Copy()
		keep = append(keep, &r)
		continue recordloop
	}
	// only return the records we keep
	pack.records = keep
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
