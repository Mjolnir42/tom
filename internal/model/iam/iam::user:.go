/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/handler"
)

// routeRegisterUser registers the library routes with the
// request router
func (m *Model) routeRegisterUser(rt *httprouter.Router) {
	rt.DELETE(`/idlib/:lib/user/:user`, m.x.Authenticated(m.UserRemove))
	rt.GET(`/idlib/:lib/user/:user`, m.x.Authenticated(m.UserShow))
	rt.GET(`/idlib/:lib/user/`, m.x.Authenticated(m.UserList))
	rt.PATCH(`/idlib/:lib/user/:user`, m.x.Authenticated(m.UserUpdate))
	rt.POST(`/idlib/:lib/user/`, m.x.Authenticated(m.UserAdd))
}

// handleRegisterUser registers user application core handlers in
// the provided handlermap
func handleRegisterUser(hm *handler.Map, length int) {
	hm.Add(NewUserReadHandler(length))
	hm.Add(NewUserWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
