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

		for m := range pack.ExportJSON(f.fFmtJSN) {
			f.outpipeJSON <- m
		}
	}

	f.outpipeIPFIX <- pack.ExportIPFIX()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
