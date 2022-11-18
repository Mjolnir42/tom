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
	proto.AssertCommandIsDefined(proto.CmdLibraryAdd)

	registry = append(registry, function{
		cmd:    proto.CmdLibraryAdd,
		handle: libraryAdd,
	})
}

func libraryAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.LibraryAdd)
}

func exportLibraryAdd(result *proto.Result, r *msg.Result) {
	result.Library = &[]proto.Library{}
	*result.Library = append(*result.Library, r.Library...)
}

// LibraryAdd function
func (m *Model) LibraryAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(
		r, params,
		proto.CmdLibraryAdd,
		msg.SectionLibrary,
		proto.ActionAdd,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Library = *req.Library

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportLibraryAdd)
}

func (h *LibraryWriteHandler) add(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
