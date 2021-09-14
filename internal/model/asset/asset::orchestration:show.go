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
	proto.AssertCommandIsDefined(proto.CmdOrchestrationShow)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationShow,
		handle: orchestrationShow,
	})
}

func orchestrationShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationShow)
}

func exportOrchestrationShow(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationShow function
func (m *Model) OrchestrationShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = proto.ActionShow
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
	m.x.Send(&w, &result, exportOrchestrationShow)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
