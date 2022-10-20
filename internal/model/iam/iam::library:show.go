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
	proto.AssertCommandIsDefined(proto.CmdLibraryShow)

	registry = append(registry, function{
		cmd:    proto.CmdLibraryShow,
		handle: libraryShow,
	})
}

func libraryShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.LibraryShow)
}

func exportLibraryShow(result *proto.Result, r *msg.Result) {
	result.Library = &[]proto.Library{}
	*result.Library = append(*result.Library, r.Library...)
}

// LibraryShow function
func (m *Model) LibraryShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdLibraryShow,
		msg.SectionLibrary,
		proto.ActionShow,
	)

	switch {
	case r.URL.Query().Get(`name`) != ``:
		request.Library.Name = r.URL.Query().Get(`name`)
	default:
		request.Library.TomID = params.ByName(`tomID`)
		if err := request.Library.ParseTomID(); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportLibraryShow)
}

// show returns full details for a specific server
func (h *LibraryReadHandler) show(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
