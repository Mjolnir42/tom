/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

// RouteRegisterServer registers the server routes with the request
// router
func (x *Rest) RouteRegisterServer(rt *httprouter.Router) *httprouter.Router {
	rt.GET(`/server/`, x.Authenticated(x.ServerList))
	rt.GET(`/server/:tomID`, x.Authenticated(x.ServerShow))
	rt.POST(`/server/`, x.Authenticated(x.ServerAdd))
	rt.DELETE(`/server/:tomID`, x.Authenticated(x.ServerRemove))
	return rt
}

// ServerList function
func (x *Rest) ServerList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// server is uniquely identified. Process this as ServerShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		x.ServerShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.hm.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerShow function
func (x *Rest) ServerShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionShow
	request.Server = proto.Server{
		TomID:     params.ByName(`tomID`),
		Namespace: r.URL.Query().Get(`namespace`),
		Name:      r.URL.Query().Get(`name`),
	}

	if err := request.Server.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			x.replyBadRequest(&w, &request, err)
			return
		}
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.hm.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerAdd function
func (x *Rest) ServerAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w, x.lm)

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionAdd

	req := proto.Server{}
	if err := decodeJSONBody(r, &req); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Server = req

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.hm.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerRemove function
func (x *Rest) ServerRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionRemove
	request.Server = proto.Server{
		TomID: params.ByName(`tomID`),
	}

	if err := request.Server.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			x.replyBadRequest(&w, &request, err)
			return
		}
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.hm.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
