/*-
 * Copyright (c) 2022, Jörg Pernfuß
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

func init() {
	proto.AssertCommandIsDefined(proto.CmdMachEnrol)

	registry = append(registry, function{
		cmd:    proto.CmdMachEnrol,
		handle: machineEnrol,
	})
}

func machineEnrol(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.MachineEnrol)
}

func exportMachineEnrol(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// MachineEnrol function
func (m *Model) MachineEnrol(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdMachEnrol,
		msg.SectionMachine,
		proto.ActionEnrolment,
	)
	request.User.LibraryName = params.ByName(`lib`)
	request.User.UserName = params.ByName(`user`)

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
	m.x.Send(&w, &result, exportMachineEnrol)
}

// update ...
func (h *UserWriteHandler) enrolment(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
