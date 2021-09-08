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
	rt.DELETE(`/idlib/:lib/team/:team`, m.x.Authenticated(m.TeamRemove))
	rt.GET(`/idlib/:lib/team/:team`, m.x.Authenticated(m.TeamShow))
	rt.GET(`/idlib/:lib/team/`, m.x.Authenticated(m.TeamList))
	rt.PATCH(`/idlib/:lib/team/:team`, m.x.Authenticated(m.TeamUpdate))
	rt.POST(`/idlib/:lib/team/`, m.x.Authenticated(m.TeamAdd))

	rt.DELETE(`/idlib/:lib/team/:team/headof`, m.x.Authenticated(m.TeamHeadOfUnset))
	rt.PUT(`/idlib/:lib/team/:team/headof`, m.x.Authenticated(m.TeamHeadOfSet))

	rt.DELETE(`/idlib/:lib/team/:team/member/:user`, m.x.Authenticated(m.TeamMemberRemove))
	rt.GET(`/idlib/:lib/team/:team/member/`, m.x.Authenticated(m.TeamMemberList))
	rt.PATCH(`/idlib/:lib/team/:team/member/`, m.x.Authenticated(m.TeamMemberAdd))
	rt.PUT(`/idlib/:lib/team/:team/member/`, m.x.Authenticated(m.TeamMemberSet))
}

// handleRegisterTeam registers team application core handlers in
// the provided handlermap
func handleRegisterTeam(hm *handler.Map, length int) {
	hm.Add(NewTeamReadHandler(length))
	hm.Add(NewTeamWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
