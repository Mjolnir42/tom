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
	proto.AssertCommandIsDefined(proto.CmdTeamList)

	registry = append(registry, function{
		cmd:    proto.CmdTeamList,
		handle: teamList,
	})
}

func teamList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.TeamList)
}

func exportTeamList(result *proto.Result, r *msg.Result) {
	result.Team = &[]proto.Team{}
	*result.Team = append(*result.Team, r.Team...)
}

// TeamList function
func (m *Model) TeamList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the namespace is uniquely
	// identified. Process this as TeamShow request
	if r.URL.Query().Get(`name`) != `` {
		m.TeamShow(w, r, params)
		return
	}

	request := msg.New(
		r, params,
		proto.CmdTeamList,
		msg.SectionTeam,
		proto.ActionList,
	)
	request.Team.LibraryName = params.ByName(`lib`)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeamList)
}

// list returns all namespaces
func (h *TeamReadHandler) list(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
