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
	proto.AssertCommandIsDefined(proto.CmdLibraryList)

	registry = append(registry, function{
		cmd:    proto.CmdLibraryList,
		handle: libraryList,
	})
}

func libraryList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.LibraryList)
}

func exportLibraryList(result *proto.Result, r *msg.Result) {
	result.Library = &[]proto.Library{}
	*result.Library = append(*result.Library, r.Library...)
}

// LibraryList function
func (m *Model) LibraryList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the library is uniquely
	// identified. Process this as LibraryShow request
	if r.URL.Query().Get(`name`) != `` {
		m.LibraryShow(w, r, params)
		return
	}

	request := msg.New(
		r, params,
		proto.CmdLibraryList,
		msg.SectionLibrary,
		proto.ActionList,
	)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportLibraryList)
}

// list returns all identity libraries
func (h *LibraryReadHandler) list(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
