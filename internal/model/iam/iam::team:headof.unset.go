/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdTeamHdUnset)

	registry = append(registry, function{
		cmd:    proto.CmdTeamHdUnset,
		handle: teamHeadOfUnset,
	})
}

func teamHeadOfUnset(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.TeamHeadOfUnset)
}

func exportTeamHeadOfUnset(result *proto.Result, r *msg.Result) {
	result.Team = &[]proto.Team{}
	*result.Team = append(*result.Team, r.Team...)
}

// TeamHeadOfUnset ...
func (m *Model) TeamHeadOfUnset(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdTeamHdUnset,
		msg.SectionTeam,
		proto.ActionHdUnset,
	)
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeamHeadOfUnset)
}

// headOfUnset ...
func (h *TeamWriteHandler) headOfUnset(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
