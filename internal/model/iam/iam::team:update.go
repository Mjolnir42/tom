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
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdTeamUpdate)

	registry = append(registry, function{
		cmd:    proto.CmdTeamUpdate,
		handle: teamUpdate,
	})
}

func teamUpdate(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.TeamUpdate)
}

func exportTeamUpdate(result *proto.Result, r *msg.Result) {
	result.Team = &[]proto.Team{}
	*result.Team = append(*result.Team, r.Team...)
}

// TeamUpdate function
func (m *Model) TeamUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdTeamUpdate,
		msg.SectionTeam,
		proto.ActionUpdate,
	)
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Update.Team = *req.Team
	request.Update.Team.LibraryName = request.Team.LibraryName

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeamUpdate)
}

// update ...
func (h *TeamWriteHandler) update(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
