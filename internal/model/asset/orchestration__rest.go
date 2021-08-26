/*-
 * Copyright (c) 2020, Jörg Pernfuß
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

// RouteRegisterOrchestration registers the runtime routes with the request
// router
func (m *Model) RouteRegisterOrchestration(rt *httprouter.Router) *httprouter.Router {
	rt.GET(`/orchestration/`, m.x.Authenticated(m.OrchestrationList))
	rt.GET(`/orchestration/:tomID`, m.x.Authenticated(m.OrchestrationShow))
	rt.POST(`/orchestration/`, m.x.Authenticated(m.OrchestrationAdd))
	rt.DELETE(`/orchestration/:tomID`, m.x.Authenticated(m.OrchestrationRemove))
	return rt
}

func exportOrchestration(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationList function
func (m *Model) OrchestrationList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// runtime is uniquely identified. Process this as OrchestrationShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.OrchestrationShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = msg.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// OrchestrationShow function
func (m *Model) OrchestrationShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = msg.ActionShow
	request.Orchestration = proto.Orchestration{
		TomID:     params.ByName(`tomID`),
		Namespace: r.URL.Query().Get(`namespace`),
		Name:      r.URL.Query().Get(`name`),
	}

	if err := request.Orchestration.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
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

// OrchestrationAdd function
func (m *Model) OrchestrationAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = msg.ActionAdd

	req := proto.Orchestration{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Orchestration = req

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result)
}

// OrchestrationRemove function
func (m *Model) OrchestrationRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = msg.ActionRemove
	request.Orchestration = proto.Orchestration{
		TomID: params.ByName(`tomID`),
	}

	if err := request.Orchestration.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
