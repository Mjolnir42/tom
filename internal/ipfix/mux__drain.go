/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
)

func (m *ipfixMux) drainDiscardChannels() {
	defer m.wg.Done()

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.discard:
			m.lm.GetLogger(`error`).
				Errorln(`ipfix.mux: drained stray message from channel discard`)
		case <-m.discardJSON:
			m.lm.GetLogger(`error`).
				Errorln(`ipfix.mux: drained stray message from channel discardJSON`)
		case <-m.discardFDR:
			m.lm.GetLogger(`error`).
				Errorln(`ipfix.mux: drained stray message from channel discardFDR`)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
