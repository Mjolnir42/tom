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

// RouteRegisterUser registers the library routes with the
// request router
func (m *Model) RouteRegisterUser(rt *httprouter.Router) *httprouter.Router {
	rt.DELETE(`/idlib/:lib/user/:user`, m.x.Authenticated(m.UserRemove))
	rt.GET(`/idlib/:lib/user/:user`, m.x.Authenticated(m.UserShow))
	rt.GET(`/idlib/:lib/user/`, m.x.Authenticated(m.UserList))
	rt.PATCH(`/idlib/:lib/user/:user`, m.x.Authenticated(m.UserUpdate))
	rt.POST(`/idlib/:lib/user/`, m.x.Authenticated(m.UserAdd))

	return rt
}

func exportUser(result *proto.Result, r *msg.Result) {
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

	request := msg.New(r, params)
	request.Section = msg.SectionUser
	request.Action = msg.ActionList
	request.User = proto.User{
		LibraryName: params.ByName(`lib`),
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// UserShow function
func (m *Model) UserShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionUser
	request.Action = msg.ActionShow
	request.User = proto.User{
		LibraryName: params.ByName(`lib`),
	}

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
	m.x.Send(&w, &result)
}

// UserAdd function
func (m *Model) UserAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionUser
	request.Action = msg.ActionAdd

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
	m.x.Send(&w, &result)
}

// UserRemove function
func (m *Model) UserRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionUser
	request.Action = msg.ActionRemove
	request.User = proto.User{
		LibraryName: params.ByName(`lib`),
		UserName:    params.ByName(`user`),
	}

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
	m.x.Send(&w, &result)
}

// UserUpdate function
func (m *Model) UserUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionUser
	request.Action = msg.ActionRemove
	request.User = proto.User{
		LibraryName: params.ByName(`lib`),
		UserName:    params.ByName(`user`),
	}

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Update.User = *req.User
	request.Update.User.LibraryName = request.User.LibraryName

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
