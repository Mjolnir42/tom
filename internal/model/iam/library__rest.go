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

func exportLibrary(result *proto.Result, r *msg.Result) {
	result.Library = &[]proto.Library{}
	*result.Library = append(*result.Library, r.Library...)
}

// LibraryList function
func (m *Model) LibraryList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the namespace is uniquely
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
	m.x.Send(&w, &result, exportLibrary)
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
		request.Library.Name = params.ByName(`lib`)
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportLibrary)
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
	m.x.Send(&w, &result, exportLibrary)
}

// LibraryRemove function
func (m *Model) LibraryRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdLibraryRemove,
		msg.SectionLibrary,
		proto.ActionRemove,
	)
	request.Library.Name = params.ByName(`lib`)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportLibrary)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
