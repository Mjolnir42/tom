/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"bytes"
	"encoding/json"
	"net"

	flow "github.com/EdgeCast/vflow/ipfix"
	"github.com/mjolnir42/flowdata"
)

func (f *procFilter) conversionWorker() {
	defer f.wg.Done()

	for {
		select {
		case <-f.quit:
			break
		case frame := <-f.pipeConvert:
			f.convert(frame)
		}
	}
}

func (f *procFilter) convert(frame IPFIXMessage) {
	// decode ipfix wireformat message
	vfMsg := f.decode(frame)
	if vfMsg == nil {
		return
	}

	// JSON marshalling of the vflow format
	vfJSON := f.marshalVFlow(vfMsg)
	if vfJSON == nil {
		return
	}

	// unfiltered JSON output in vflow format was requested
	if f.fOutJSN && f.fRawJSN && f.fFmtJSN == `vflow` {
		j := make([]byte, len(vfJSON))
		copy(j, vfJSON)
		select {
		case f.outpipeRawJS <- j:
		default:
		}
	}

	// JSON unmarshal into flowdata.Message format
	fdMsg := &flowdata.Message{}
	if err := json.Unmarshal(vfJSON, fdMsg); err != nil {
		f.err <- err
		return
	}

	pack := MessagePack{
		raddr: frame.raddr,
		header: IPFIXHeader{
			Version:    vfMsg.Header.Version,
			Length:     vfMsg.Header.Length,
			ExportTime: vfMsg.Header.ExportTime,
			SequenceNo: vfMsg.Header.SequenceNo,
			DomainID:   vfMsg.Header.DomainID,
		},
		records: make([]*flowdata.Record, len(vfMsg.DataSets)),
		jsons:   make([][]byte, len(vfMsg.DataSets)),
	}
	pack.SetClientID()

recordloop:
	for r := range fdMsg.Convert() {
		if r.SrcAddress != `` {
			p := net.ParseIP(r.SrcAddress)
			if p == nil {
				continue recordloop
			}
			r.SrcAddress = p.String()
			if r.DstAddress != `` {
				p := net.ParseIP(r.DstAddress)
				if p == nil {
					continue recordloop
				}
				r.DstAddress = p.String()
			}
			if r.AgentID != `` {
				p := net.ParseIP(r.AgentID)
				if p == nil {
					continue recordloop
				}
				r.AgentID = p.String()
			}
			// unfiltered JSON output in vflow format was requested
			if f.fOutJSN && f.fRawJSN && f.fFmtJSN == `flowdata` {
				j, err := json.Marshal(&OutputRecord{
					OctetCount:  r.OctetCount,
					PacketCount: r.PacketCount,
					ProtocolID:  r.ProtocolID,
					Protocol:    r.Protocol,
					IPVersion:   r.IPVersion,
					SrcAddress:  r.SrcAddress,
					SrcPort:     r.SrcPort,
					DstAddress:  r.DstAddress,
					DstPort:     r.DstPort,
					TcpFlags:    r.TcpFlags.Copy(),
					StartMilli:  r.StartMilli,
					EndMilli:    r.EndMilli,
					AgentID:     r.AgentID,
				})
				if err != nil {
					f.err <- err
					return
				}
				select {
				case f.outpipeRawJS <- j:
				default:
				}
			}

			// store record
			pack.records = append(pack.records, &r)
		}
		f.pipeFilter <- pack
	}
}

func (f *procFilter) decode(frame IPFIXMessage) *flow.Message {
	decoder := flow.NewDecoder(*frame.raddr, frame.body)

	decodedMsg, err := decoder.Decode(mCache)
	if err != nil {
		f.err <- err
		return nil
	}
	if decodedMsg == nil {
		return nil
	}
	if len(decodedMsg.DataSets) == 0 {
		return nil
	}
	return decodedMsg
}

func (f *procFilter) marshalVFlow(m *flow.Message) []byte {
	buf := new(bytes.Buffer)
	b, err := m.JSONMarshal(buf)
	if err != nil {
		f.err <- err
		return nil
	}
	return b
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
