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
	proto.AssertCommandIsDefined(proto.CmdUserRemove)

	registry = append(registry, function{
		cmd:    proto.CmdUserRemove,
		handle: userRemove,
	})
}

func userRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.UserRemove)
}

func exportUserRemove(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// UserRemove function
func (m *Model) UserRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdUserRemove,
		msg.SectionUser,
		proto.ActionRemove,
	)
	request.User.LibraryName = params.ByName(`lib`)
	request.User.UserName = params.ByName(`user`)

	if err := request.Namespace.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportUserRemove)
}

// remove ...
func (h *UserWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
