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
	proto.AssertCommandIsDefined(proto.CmdUserShow)

	registry = append(registry, function{
		cmd:    proto.CmdUserShow,
		handle: userShow,
	})
}

func userShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.UserShow)
}

func exportUserShow(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// UserShow function
func (m *Model) UserShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdUserShow,
		msg.SectionUser,
		proto.ActionShow,
	)
	request.User.LibraryName = params.ByName(`lib`)

	switch {
	case r.URL.Query().Get(`name`) != ``:
		request.User.UserName = r.URL.Query().Get(`name`)
	default:
		request.User.UserName = params.ByName(`user`)
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportUserShow)
}

// show returns full details for a specific server
func (h *UserReadHandler) show(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
