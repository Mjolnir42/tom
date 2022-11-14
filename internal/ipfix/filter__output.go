/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
)

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
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
