/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

// routeRegisterServer registers the server routes with the request
// router
func (m *Model) routeRegisterServer(rt *httprouter.Router) {
	rt.GET(`/server/`, m.x.Authenticated(m.ServerList))
	rt.GET(`/server/:tomID`, m.x.Authenticated(m.ServerShow))
	rt.POST(`/server/`, m.x.Authenticated(m.ServerAdd))
	rt.DELETE(`/server/:tomID`, m.x.Authenticated(m.ServerRemove))
}

func exportServer(result *proto.Result, r *msg.Result) {
	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)
}

// ServerList function
func (m *Model) ServerList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// server is uniquely identified. Process this as ServerShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.ServerShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportServer)
}

// ServerShow function
func (m *Model) ServerShow(w http.ResponseWriter, r *http.Request,
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
		if !(err == proto.ErrEmptyTomID && request.Server.Name != ``) {
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
	m.x.Send(&w, &result, exportServer)
}

// ServerAdd function
func (m *Model) ServerAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Server = *req.Server

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportServer)
}

// ServerRemove function
func (m *Model) ServerRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionRemove
	request.Server = proto.Server{
		TomID: params.ByName(`tomID`),
	}

	if err := request.Server.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportServer)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
