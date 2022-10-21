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
	proto.AssertCommandIsDefined(proto.CmdUserAdd)

	registry = append(registry, function{
		cmd:    proto.CmdUserAdd,
		handle: userAdd,
	})
}

func userAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.UserAdd)
}

func exportUserAdd(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// UserAdd function
func (m *Model) UserAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(
		r, params,
		proto.CmdUserAdd,
		msg.SectionUser,
		proto.ActionAdd,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	req.User.LibraryName = params.ByName(`lib`)
	request.User = *req.User

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportUserAdd)
}

// add ...
func (h *UserWriteHandler) add(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
