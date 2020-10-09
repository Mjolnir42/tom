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

func (x *Rest) RouteRegisterServer(rt *httprouter.Router) *httprouter.Router {
	rt.GET(`/server/`, x.Authenticated(x.ServerList))
	rt.GET(`/server/:serverID`, x.Authenticated(x.ServerShow))
	return rt
}

// ServerList function
func (x *Rest) ServerList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// server is uniquely identified. Process this as ServerShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		return x.ServerShow(w, r, params)
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
		ID:        params.ByName(`serverID`),
		Namespace: r.URL.Query().Get(`namespace`),
		Name:      r.URL.Query().Get(`name`),
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
