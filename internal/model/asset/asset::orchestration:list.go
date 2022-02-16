/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"database/sql"
	"net/http"
	"time"

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

	request := msg.New(
		r, params,
		proto.CmdOrchestrationList,
		msg.SectionOrchestration,
		proto.ActionList,
	)

	if r.URL.Query().Get(`namespace`) != `` {
		request.Orchestration.Namespace = r.URL.Query().Get(`namespace`)
		if err := proto.ValidNamespace(request.Orchestration.Namespace); err != nil {
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
	m.x.Send(&w, &result, exportOrchestrationList)
}

// list returns all orchestration environments
func (h *OrchestrationReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		id, dictName, key, value, author string
		creationTime                     time.Time
		rows                             *sql.Rows
		err                              error
		namespace                        sql.NullString
		orch                             proto.OrchestrationHeader
		ok                               bool
	)

	if q.Orchestration.Namespace != `` {
		namespace.String = q.Orchestration.Namespace
		namespace.Valid = true
	}

	list := make(map[string]proto.OrchestrationHeader)
	if rows, err = h.stmtList.Query(
		namespace,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&dictName,
			&key,
			&value,
			&author,
			&creationTime,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		if orch, ok = list[id]; !ok {
			orch = proto.OrchestrationHeader{}
		}
		orch.Namespace = dictName
		switch key {
		case `type`:
			orch.Type = value
		case `name`:
			orch.Name = value
			orch.CreatedBy = author
			orch.CreatedAt = creationTime.Format(msg.RFC3339Milli)
		}
		list[id] = orch
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	for _, oenv := range list {
		mr.OrchestrationHeader = append(mr.OrchestrationHeader, oenv)
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
