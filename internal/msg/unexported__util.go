/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

// requestID extracts the RequestID set by Basic Authentication, making
// the ID consistent between all logs
func requestID(params httprouter.Params) (id uuid.UUID) {
	id, _ = uuid.FromString(params.ByName(`RequestID`))
	return
}

// authUser extracts the AuthenticatedUser set by Basic Authentication
func authUser(params httprouter.Params) []string {
	return strings.Split(params.ByName(`AuthenticatedUser`), `~`)
}

// authEnforcement extracts the information if auth enforcement is
// active
func authEnforcement(params httprouter.Params) bool {
	switch params.ByName(`AuthenticationEnforcement`) {
	case `true`:
		return true
	case `false`:
		return false
	default:
		return false
	}
}

// remoteAddr extracts the IP address part of the IP:port string
// set as net/http.Request.RemoteAddr. It handles IPv4 cases like
// 192.0.2.1:48467 and IPv6 cases like [2001:db8::1%lo0]:48467
func remoteAddr(r *http.Request) string {
	var addr string

	switch {
	case strings.Contains(r.RemoteAddr, `]`):
		// IPv6 address [2001:db8::1%lo0]:48467
		addr = strings.Split(r.RemoteAddr, `]`)[0]
		addr = strings.Split(addr, `%`)[0]
		addr = strings.TrimLeft(addr, `[`)
	default:
		// IPv4 address 192.0.2.1:48467
		addr = strings.Split(r.RemoteAddr, `:`)[0]
	}
	return addr
}

// requestURI extracts the recorded RequestURI set by Basic
// Authentication
func requestURI(params httprouter.Params) string {
	return params.ByName(`RequestURI`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
