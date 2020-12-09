/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

// RouteRegisterNamespace registers the namespace routes with the
// request router
func (m *Model) RouteRegisterNamespace(rt *httprouter.Router) *httprouter.Router {
	rt.GET(`/namespace/`, m.x.Authenticated(m.NamespaceList))
	rt.GET(`/namespace/:tomID`, m.x.Authenticated(m.NamespaceShow))
	rt.POST(`/namespace/`, m.x.Authenticated(m.NamespaceAdd))
	rt.DELETE(`/namespace/:tomID`, m.x.Authenticated(m.NamespaceRemove))

	rt.POST(`/namespace/:tomID/attribute/`, m.x.Authenticated(m.NamespaceAttributeAdd))

	rt.PUT(`/namespace/:tomID/property/`, m.x.Authenticated(m.NamespacePropertySet))
	rt.PATCH(`/namespace/:tomID/property/`, m.x.Authenticated(m.NamespacePropertyUpdate))
	return rt
}

// NamespaceList function
func (m *Model) NamespaceList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the namespace is uniquely
	// identified. Process this as NamespaceShow request
	if r.URL.Query().Get(`name`) != `` {
		m.NamespaceShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// NamespaceShow function
func (m *Model) NamespaceShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionShow
	request.Namespace = proto.Namespace{
		TomID: params.ByName(`tomID`),
		Name:  r.URL.Query().Get(`name`),
	}

	if err := request.Namespace.ParseTomID(); err != nil {
		if !(err == proto.ErrEmptyTomID && request.Namespace.Name != ``) {
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
	m.x.Send(&w, &result)
}

// NamespaceAdd function
func (m *Model) NamespaceAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionAdd

	req := proto.Namespace{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = req

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// NamespaceRemove function
func (m *Model) NamespaceRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionRemove
	request.Namespace = proto.Namespace{
		TomID: params.ByName(`tomID`),
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

// NamespaceAttributeAdd function
func (m *Model) NamespaceAttributeAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionAttrAdd

	req := proto.Namespace{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = req

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// NamespacePropertySet function
func (m *Model) NamespacePropertySet(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionPropSet

	req := proto.Namespace{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = req

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// NamespacePropertyUpdate function
func (m *Model) NamespacePropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionPropUpdate

	req := proto.Namespace{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = req

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
