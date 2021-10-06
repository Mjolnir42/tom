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
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdServerRemove)

	registry = append(registry, function{
		cmd:    proto.CmdServerRemove,
		handle: serverRemove,
	})
}

func serverRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ServerRemove)
}

func exportServerRemove(result *proto.Result, r *msg.Result) {
	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)
}

// ServerRemove function
func (m *Model) ServerRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = proto.ActionRemove
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
	m.x.Send(&w, &result, exportServerRemove)
}

// remove invalidates an existing server
func (h *ServerWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
