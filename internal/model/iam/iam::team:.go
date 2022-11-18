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

// routeRegisterTeam registers the team routes with the
// request router
func (m *Model) routeRegisterTeam(rt *httprouter.Router) {
}

// handleRegisterTeam registers team application core handlers in
// the provided handlermap
func handleRegisterTeam(hm *handler.Map, length int) {
	hm.Add(NewTeamReadHandler(length))
	hm.Add(NewTeamWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
