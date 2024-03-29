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
	result.RuntimeHeader = &[]proto.RuntimeHeader{}
	*result.RuntimeHeader = append(*result.RuntimeHeader, r.RuntimeHeader...)
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

	request := msg.New(
		r, params,
		proto.CmdRuntimeList,
		msg.SectionRuntime,
		proto.ActionList,
	)

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}
	if r.URL.Query().Get(`namespace`) != `` {
		request.Runtime.Namespace = r.URL.Query().Get(`namespace`)
		if err := proto.ValidNamespace(request.Runtime.Namespace); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportRuntimeList)
}

// list returns all servers
func (h *RuntimeReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		rteID, rteNs, key, value, author string
		creationTime                     time.Time
		namespace                        sql.NullString
		rows                             *sql.Rows
		rte                              proto.RuntimeHeader
		err                              error
		ok                               bool
	)

	if q.Runtime.Namespace != `` {
		namespace.String = q.Runtime.Namespace
		namespace.Valid = true
	}

	list := make(map[string]proto.RuntimeHeader)
	if rows, err = h.stmtList.Query(
		namespace,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&rteID,
			&rteNs,
			&key,
			&value,
			&author,
			&creationTime,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		if rte, ok = list[rteID]; !ok {
			rte = proto.RuntimeHeader{}
		}
		rte.Namespace = rteNs
		switch key {
		case `type`:
			rte.Type = value
		case `name`:
			rte.Name = value
			rte.CreatedBy = author
			rte.CreatedAt = creationTime.Format(msg.RFC3339Milli)
		}
		list[rteID] = rte
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	for _, rte := range list {
		mr.RuntimeHeader = append(mr.RuntimeHeader, rte)
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
