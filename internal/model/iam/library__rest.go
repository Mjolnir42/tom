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

// RouteRegisterLibrary registers the library routes with the
// request router
func (m *Model) RouteRegisterLibrary(rt *httprouter.Router) *httprouter.Router {
	rt.GET(`/idlib/`, m.x.Authenticated(m.LibraryList))
	rt.GET(`/idlib/:lib`, m.x.Authenticated(m.LibraryShow))
	rt.POST(`/idlib/`, m.x.Authenticated(m.LibraryAdd))
	rt.DELETE(`/idlib/:lib`, m.x.Authenticated(m.LibraryRemove))

	return rt
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

	request := msg.New(r, params)
	request.Section = msg.SectionLibrary
	request.Action = msg.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// LibraryShow function
func (m *Model) LibraryShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionLibrary
	request.Action = msg.ActionShow
	request.Library = proto.Library{}

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
	m.x.Send(&w, &result)
}

// LibraryAdd function
func (m *Model) LibraryAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionLibrary
	request.Action = msg.ActionAdd

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
	m.x.Send(&w, &result)
}

// LibraryRemove function
func (m *Model) LibraryRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionLibrary
	request.Action = msg.ActionRemove
	request.Library = proto.Library{
		Name: params.ByName(`lib`),
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
