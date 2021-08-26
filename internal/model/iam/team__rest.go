/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

// routeRegisterTeam registers the team routes with the
// request router
func (m *Model) routeRegisterTeam(rt *httprouter.Router) {
	rt.DELETE(`/idlib/:lib/team/:team`, m.x.Authenticated(m.TeamRemove))
	rt.GET(`/idlib/:lib/team/:team`, m.x.Authenticated(m.TeamShow))
	rt.GET(`/idlib/:lib/team/`, m.x.Authenticated(m.TeamList))
	rt.PATCH(`/idlib/:lib/team/:team`, m.x.Authenticated(m.TeamUpdate))
	rt.POST(`/idlib/:lib/team/`, m.x.Authenticated(m.TeamAdd))

	rt.DELETE(`/idlib/:lib/team/:team/headof`, m.x.Authenticated(m.TeamHeadOfUnset))
	rt.PUT(`/idlib/:lib/team/:team/headof`, m.x.Authenticated(m.TeamHeadOfSet))

	rt.DELETE(`/idlib/:lib/team/:team/member/:user`, m.x.Authenticated(m.TeamMemberRemove))
	rt.GET(`/idlib/:lib/team/:team/member/`, m.x.Authenticated(m.TeamMemberList))
	rt.PATCH(`/idlib/:lib/team/:team/member/`, m.x.Authenticated(m.TeamMemberAdd))
	rt.PUT(`/idlib/:lib/team/:team/member/`, m.x.Authenticated(m.TeamMemberSet))
}

func exportTeam(result *proto.Result, r *msg.Result) {
	result.Team = &[]proto.Team{}
	*result.Team = append(*result.Team, r.Team...)
}

// TeamList function
func (m *Model) TeamList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the namespace is uniquely
	// identified. Process this as TeamShow request
	if r.URL.Query().Get(`name`) != `` {
		m.TeamShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionList
	request.Team = proto.Team{
		LibraryName: params.ByName(`lib`),
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamShow function
func (m *Model) TeamShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionShow
	request.Team = proto.Team{
		LibraryName: params.ByName(`lib`),
	}

	switch {
	case r.URL.Query().Get(`name`) != ``:
		request.Team.TeamName = r.URL.Query().Get(`name`)
	default:
		request.Team.TeamName = params.ByName(`team`)
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamAdd function
func (m *Model) TeamAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Team = *req.Team
	request.Team.LibraryName = params.ByName(`lib`)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamRemove function
func (m *Model) TeamRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionRemove
	request.Team = proto.Team{
		LibraryName: params.ByName(`lib`),
		TeamName:    params.ByName(`team`),
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamUpdate function
func (m *Model) TeamUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionUpdate
	request.Team = proto.Team{
		LibraryName: params.ByName(`lib`),
		TeamName:    params.ByName(`team`),
	}

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Update.Team = *req.Team
	request.Update.Team.LibraryName = request.Team.LibraryName

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamHeadOfSet ...
func (m *Model) TeamHeadOfSet(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionHdSet

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Team = *req.Team
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	request.Team.TeamLead.LibraryName = request.Team.LibraryName
	if request.Team.TeamLead.UserName == `` {
		m.x.ReplyBadRequest(&w, &request, fmt.Errorf("Found team-lead with empty username"))
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamHeadOfUnset ...
func (m *Model) TeamHeadOfUnset(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionHdUnset
	request.Team = proto.Team{
		LibraryName: params.ByName(`lib`),
		TeamName:    params.ByName(`team`),
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamMemberAdd ...
func (m *Model) TeamMemberAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionMbrAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Team = *req.Team
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	for i := range *(request.Team.Member) {
		// stamp all referenced users with the team's LibraryName
		(*request.Team.Member)[i].LibraryName = request.Team.LibraryName
		// check that the username is not empty
		if (*request.Team.Member)[i].UserName == `` {
			m.x.ReplyBadRequest(&w, &request, fmt.Errorf("Found member with empty username"))
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamMemberSet ...
func (m *Model) TeamMemberSet(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionMbrSet

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Team = *req.Team
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	for i := range *(request.Team.Member) {
		// stamp all referenced users with the team's LibraryName
		(*request.Team.Member)[i].LibraryName = request.Team.LibraryName
		// check that the username is not empty
		if (*request.Team.Member)[i].UserName == `` {
			m.x.ReplyBadRequest(&w, &request, fmt.Errorf("Found member with empty username"))
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamMemberList ...
func (m *Model) TeamMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionMbrList
	request.Team = proto.Team{
		LibraryName: params.ByName(`lib`),
		TeamName:    params.ByName(`team`),
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// TeamMemberRemove ...
func (m *Model) TeamMemberRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionMbrRemove

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Team = *req.Team
	request.Team.LibraryName = params.ByName(`lib`)
	request.Team.TeamName = params.ByName(`team`)

	for i := range *(request.Team.Member) {
		// stamp all referenced users with the team's LibraryName
		(*request.Team.Member)[i].LibraryName = request.Team.LibraryName
		// check that the username is not empty
		if (*request.Team.Member)[i].UserName == `` {
			m.x.ReplyBadRequest(&w, &request, fmt.Errorf("Found member with empty username"))
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportTeam)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
