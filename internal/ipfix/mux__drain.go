/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"fmt"
)

func (m *ipfixMux) drainDiscardChannels() {
	defer m.wg.Done()

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.exit:
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

// opportunUDP is an opportunistic drain for channel outUDP that reads
// from the channel if no reader is configured.
func (m *ipfixMux) opportunUDP() {
	defer m.wg.Done()
	if m.fOutUDP {
		// udpClient is active
		return
	}

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.exit:
			break drainloop
		case <-m.outUDP:
			m.err <- fmt.Errorf(
				"ipfix.mux: drained stray message from channel %s", `outUDP`,
			)
		}
	}
}

// opportunTCP is an opportunistic drain for channel outTCP that reads
// from the channel if no reader is configured.
func (m *ipfixMux) opportunTCP() {
	defer m.wg.Done()
	if m.fOutTCP {
		// tcpClient is active
		return
	}

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.exit:
			break drainloop
		case <-m.outTCP:
			m.err <- fmt.Errorf(
				"ipfix.mux: drained stray message from channel %s", `outTCP`,
			)
		}
	}
}

// opportunTLS is an opportunistic drain for channel outTLS that reads
// from the channel if no reader is configured.
func (m *ipfixMux) opportunTLS() {
	defer m.wg.Done()
	if m.fOutTLS {
		// tlsClient is active
		return
	}

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.exit:
			break drainloop
		case <-m.outTLS:
			m.err <- fmt.Errorf(
				"ipfix.mux: drained stray message from channel %s", `outTLS`,
			)
		}
	}
}

// opportunJSON is an opportunistic drain for channel outJSON that reads
// from the channel if no reader is configured.
func (m *ipfixMux) opportunJSON() {
	defer m.wg.Done()
	if m.fOutJSN {
		// jsonClient is active
		return
	}

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.exit:
			break drainloop
		case <-m.outJSN:
			m.err <- fmt.Errorf(
				"ipfix.mux: drained stray message from channel %s", `outJSN`,
			)
		}
	}
}

// opportunAGG is an opportunistic drain for channel outAGG that reads
// from the channel if no reader is configured.
func (m *ipfixMux) opportunAGG() {
	defer m.wg.Done()
	if m.aggregation {
		// aggregation is active
		return
	}

drainloop:
	for {
		select {
		case <-m.quit:
			break drainloop
		case <-m.exit:
			break drainloop
		case <-m.outAGG:
			m.err <- fmt.Errorf(
				"ipfix.mux: drained stray message from channel %s", `outAGG`,
			)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
