/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
// "github.com/mjolnir42/flowdata"
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

func (f *procFilter) filter(mp MessagePack) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
