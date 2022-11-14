/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"fmt"

	"github.com/mjolnir42/flowdata"
)

// setup is the function that wires the channel routing for m
func (m *ipfixMux) setup() {
	defer m.wg.Done()

	// check for invalid configurations
	if m.fInJSN && !m.fOutJSN {
		m.err <- fmt.Errorf("%s",
			`JSON input can only be forwarded to JSON output`,
		)
		close(m.exit)
		return
	}
	m.wg.Add(1)
	go m.connectJSONChannel()

	if !m.forwarding && !m.aggregation {
		m.err <- fmt.Errorf("%s",
			`No IPFIX output module or aggregation activated`,
		)
		close(m.exit)
		return
	}

	// start opportunistic drain of deactivated outputs
	m.wg.Add(5)
	go m.opportunUDP()
	go m.opportunTCP()
	go m.opportunTLS()
	go m.opportunJSON()
	go m.opportunAGG()

	// if filtering is disabled, nobody is reading from outFLT and
	// writing back into inFLT.
	// the filtering module is also started if there is a JSON output
	if !m.filtering && !m.fOutJSN {
		m.wg.Add(1)
		go m.connectFilterChannel()
	}

}

// connectFilterChannel is a direct connection between the IPFIX mux input
// and output for the filter processor. It is enabled if no filter module
// is activated
func (m *ipfixMux) connectFilterChannel() {
	defer m.wg.Done()

	for {
		select {
		case <-m.quit:
			break
		case <-m.exit:
			break
		case frame := <-m.outFLT:
			select {
			case m.inFLT <- frame:
			default:
			}
		}
	}
}

// connectJSONChannel is a direct connection between the JSON input and
// output as received JSON data is not further processed
func (m *ipfixMux) connectJSONChannel() {
	defer m.wg.Done()

	for {
		select {
		case <-m.quit:
			break
		case <-m.exit:
			break
		case buf := <-m.inJSN:
			select {
			case m.outJSN <- buf:
			default:
			}
		}
	}
}

// pipe returns the requested channel p from m
func (m *ipfixMux) pipe(p string) chan IPFIXMessage {
	switch p {
	case `inUDP`:
		return m.inUDP
	case `inTCP`:
		return m.inTCP
	case `inTLS`:
		return m.inTLS
	case `outUDP`:
		return m.outUDP
	case `outTCP`:
		return m.outTCP
	case `outTLS`:
		return m.outTLS
	case `outFLT`:
		return m.outFLT
	case `inFLT`:
		return m.inFLT
	default:
		return m.discard
	}
}

// jsonPipe returns the requested channel p from m
func (m *ipfixMux) jsonPipe(p string) chan []byte {
	switch p {
	case `inJSN`:
		return m.inJSN
	case `outJSN`:
		return m.outJSN
	case `inFLJ`:
		return m.inFLJ
	default:
		return m.discardJSON
	}
}

// flowdataPipe returns the requested channel p from m
func (m *ipfixMux) flowdataPipe(p string) chan flowdata.Record {
	switch p {
	case `inFLR`:
		return m.inFLR
	case `outAGG`:
		return m.outAGG
	default:
		return m.discardFDR
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
