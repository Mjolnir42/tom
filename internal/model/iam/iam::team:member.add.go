/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdTeamMbrAdd)

	registry = append(registry, function{
		cmd:    proto.CmdTeamMbrAdd,
		handle: teamMemberAdd,
	})
}

func teamMemberAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.TeamMemberAdd)
}

func exportTeamMemberAdd(result *proto.Result, r *msg.Result) {
	result.Team = &[]proto.Team{}
	*result.Team = append(*result.Team, r.Team...)
}

// TeamMemberAdd ...
func (m *Model) TeamMemberAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdTeamMbrAdd,
		msg.SectionTeam,
		proto.ActionMbrAdd,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Team = *req.Team
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	for i := range *(request.Team.Member) {
		// stamp all referenced users with the team's LibraryName
		(*request.Team.Member)[i].LibraryName = request.Team.LibraryName
		// check that the username is not empty
		if (*request.Team.Member)[i].UserName == `` {
			m.x.ReplyBadRequest(&w, &request, fmt.Errorf("Found member with empty username"))
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeamMemberAdd)
}

// memberAdd ...
func (h *TeamWriteHandler) memberAdd(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
