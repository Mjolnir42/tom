/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/handler"
)

// routeRegisterServer registers the server routes with the request
// router
func (m *Model) routeRegisterServer(rt *httprouter.Router) {
	rt.GET(`/server/`, m.x.Authenticated(m.ServerList))
	rt.GET(`/server/:tomID`, m.x.Authenticated(m.ServerShow))
	rt.POST(`/server/`, m.x.Authenticated(m.ServerAdd))
	rt.DELETE(`/server/:tomID`, m.x.Authenticated(m.ServerRemove))
}

// handleRegisterServer registers the server application core handlers
// in the provided handlermap
func handleRegisterServer(hm *handler.Map, length int) {
	hm.Add(NewServerReadHandler(length))
	hm.Add(NewServerWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
