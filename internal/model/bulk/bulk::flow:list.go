/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/model/bulk/"

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdFlowList)

	registry = append(registry, function{
		cmd:    proto.CmdFlowList,
		handle: flowList,
	})
}

func flowList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.FlowList)
}

func exportFlowList(result *proto.Result, r *msg.Result) {
	result.Flow = &[]proto.Flow{}
	*result.Flow = append(*result.Flow, r.Flow...)
}

// FlowList function
func (m *Model) FlowList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// flow is uniquely identified. Process this as FlowShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.FlowShow(w, r, params)
		return
	}

	request := msg.New(
		r, params,
		proto.CmdFlowList,
		msg.SectionFlow,
		proto.ActionList,
	)

	if r.URL.Query().Get(`namespace`) != `` {
		request.Flow.Namespace = r.URL.Query().Get(`namespace`)
		if err := proto.ValidNamespace(request.Flow.Namespace); err != nil {
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
	m.x.Send(&w, &result, exportFlowList)
}

// list returns all flows
func (h *FlowReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		flowID, nsName, key, value, author string
		creationTime                       time.Time
		namespace                          sql.NullString
		rows                               *sql.Rows
		flow                               proto.Flow
		err                                error
		ok                                 bool
	)

	if q.Flow.Namespace != `` {
		namespace.String = q.Flow.Namespace
		namespace.Valid = true
	}

	list := make(map[string]proto.Flow)
	if rows, err = h.stmtList.Query(
		namespace,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&flowID,
			&nsName,
			&key,
			&value,
			&author,
			&creationTime,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		if flow, ok = list[flowID]; !ok {
			flow = *(proto.NewFlow())
		}
		flow.Namespace = nsName
		switch {
		case key == `type`:
			flow.Type = value
		case key == `name`:
			flow.Name = value
			flow.CreatedBy = author
			flow.CreatedAt = creationTime.Format(msg.RFC3339Milli)
		}
		list[flowID] = flow
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	for _, flow := range list {
		mr.Flow = append(mr.Flow, flow)
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
