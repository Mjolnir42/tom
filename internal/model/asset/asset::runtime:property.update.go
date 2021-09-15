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
	proto.AssertCommandIsDefined(proto.CmdRuntimePropUpdate)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimePropUpdate,
		handle: runtimePropUpdate,
	})
}

func runtimePropUpdate(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimePropUpdate)
}

func exportRuntimePropUpdate(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimePropUpdate function
func (m *Model) RuntimePropUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionRuntime
	request.Action = proto.ActionPropUpdate
	request.Runtime = proto.Runtime{
		TomID: params.ByName(`tomID`),
	}

	if err := request.Runtime.ParseTomID(); err != nil {
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
	m.x.Send(&w, &result, exportRuntimePropUpdate)
}

// propertyUpdate ...
func (h *RuntimeWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	// TODO
	mr.NotImplemented()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
