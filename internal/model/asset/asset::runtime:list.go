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
	proto.AssertCommandIsDefined(proto.CmdRuntimeList)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeList,
		handle: runtimeList,
	})
}

func runtimeList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeList)
}

func exportRuntimeList(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimeList function
func (m *Model) RuntimeList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// runtime is uniquely identified. Process this as RuntimeShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.RuntimeShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionRuntime
	request.Action = proto.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportRuntimeList)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
