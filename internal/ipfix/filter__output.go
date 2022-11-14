/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"encoding/json"
	"time"

	"github.com/mjolnir42/flowdata"
)

type OutputRecord struct {
	OctetCount  uint64         `json:"OctetCount"`
	PacketCount uint64         `json:"PacketCount"`
	ProtocolID  uint8          `json:"ProtocolID"`
	Protocol    string         `json:"Protocol,omitempty"`
	IPVersion   uint8          `json:"IPVersion"`
	SrcAddress  string         `json:"SrcAddress"`
	SrcPort     uint16         `json:"SrcPort"`
	DstAddress  string         `json:"DstAddress"`
	DstPort     uint16         `json:"DstPort"`
	TcpFlags    flowdata.Flags `json:"TcpFlags"`
	StartMilli  time.Time      `json:"StartDateTimeMilli"`
	EndMilli    time.Time      `json:"EndDateTimeMilli"`
	AgentID     string         `json:"AgentID"`
}

func (f *procFilter) outputWorker() {
	defer f.wg.Done()

loop:
	for {
		select {
		case <-f.quit:
			break loop
		case <-f.exit:
			break loop
		case mp := <-f.pipeOutput:
			f.output(mp)
		}
	}
}

func (f *procFilter) output(pack MessagePack) {
	if f.fOutJSN && !f.fRawJSN {
		f.formatOutputJSON(&pack, f.fFmtJSN)

		for m := range pack.ExportJSON(f.fFmtJSN) {
			f.outpipeJSON <- m
		}
	}

	f.outpipeIPFIX <- pack.ExportIPFIX()

	for _, r := range pack.records {
		f.outpipeFDR <- *r
	}
}

func (f *procFilter) formatOutputJSON(pack *MessagePack, s string) {
	switch s {
	case `vflow`:
		// decode ipfix wireformat message
		vfMsg := f.decode(pack.ExportIPFIX())
		if vfMsg == nil {
			pack.jsons[0] = nil
			return
		}
		// JSON marshalling of the vflow format
		vfJSON := f.marshalVFlow(vfMsg)
		if vfJSON == nil {
			pack.jsons[0] = nil
			return
		}
		pack.jsons[0] = vfJSON
		return

	case `flowdata`:
		converr := make([]int, 0)
		for i := range pack.records {
			var err error
			pack.jsons[i], err = json.Marshal(&OutputRecord{
				OctetCount:  pack.records[i].OctetCount,
				PacketCount: pack.records[i].PacketCount,
				ProtocolID:  pack.records[i].ProtocolID,
				Protocol:    pack.records[i].Protocol,
				IPVersion:   pack.records[i].IPVersion,
				SrcAddress:  pack.records[i].SrcAddress,
				SrcPort:     pack.records[i].SrcPort,
				DstAddress:  pack.records[i].DstAddress,
				DstPort:     pack.records[i].DstPort,
				TcpFlags:    pack.records[i].TcpFlags.Copy(),
				StartMilli:  pack.records[i].StartMilli,
				EndMilli:    pack.records[i].EndMilli,
				AgentID:     pack.records[i].AgentID,
			})
			if err != nil {
				converr = append(converr, i)
				f.err <- err
			}
		}
		for i := range converr {
			x := converr[i]
			pack.jsons = append(
				pack.jsons[:x],
				pack.jsons[x+1:]...,
			)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
