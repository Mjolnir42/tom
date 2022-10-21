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
	proto.AssertCommandIsDefined(proto.CmdUserList)

	registry = append(registry, function{
		cmd:    proto.CmdUserList,
		handle: userList,
	})
}

func userList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.UserList)
}

func exportUserList(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// UserList function
func (m *Model) UserList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the namespace is uniquely
	// identified. Process this as UserShow request
	if r.URL.Query().Get(`name`) != `` {
		m.UserShow(w, r, params)
		return
	}

	request := msg.New(
		r, params,
		proto.CmdUserList,
		msg.SectionUser,
		proto.ActionList,
	)
	request.User.LibraryName = params.ByName(`lib`)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportUserList)
}

// list returns all namespaces
func (h *UserReadHandler) list(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
