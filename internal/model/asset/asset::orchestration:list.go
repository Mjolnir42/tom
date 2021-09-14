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
	proto.AssertCommandIsDefined(proto.CmdOrchestrationList)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationList,
		handle: orchestrationList,
	})
}

func orchestrationList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationList)
}

func exportOrchestrationList(result *proto.Result, r *msg.Result) {
	result.OrchestrationHeader = &[]proto.OrchestrationHeader{}
	*result.OrchestrationHeader = append(*result.OrchestrationHeader, r.OrchestrationHeader...)
}

// OrchestrationList function
func (m *Model) OrchestrationList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// orchestration is uniquely identified. Process this as OrchestrationShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.OrchestrationShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = proto.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportOrchestrationList)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
