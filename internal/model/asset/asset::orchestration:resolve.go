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
	proto.AssertCommandIsDefined(proto.CmdOrchestrationResolve)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationResolve,
		handle: orchestrationResolve,
	})
}

func orchestrationResolve(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationResolve)
}

func exportOrchestrationResolve(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	result.ServerHeader = &[]proto.ServerHeader{}
	*result.ServerHeader = append(*result.ServerHeader, r.ServerHeader...)
}

// OrchestrationResolve function
func (m *Model) OrchestrationResolve(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = proto.ActionResolve
	request.Orchestration = proto.Orchestration{
		TomID: params.ByName(`tomID`),
		Type:  params.ByName(`level`), // resolution detail type
	}

	if err := request.Orchestration.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportOrchestrationResolve)
}

// resolve ...
func (h *OrchestrationReadHandler) resolve(q *msg.Request, mr *msg.Result) {
	var (
		nsName, nsID, oreID, srvName, srvType  string
		createdAt, createdBy, namedAt, namedBy string
		since, until                           time.Time
		rows                                   *sql.Rows
		err                                    error
	)

	if err = h.stmtTxShow.QueryRow(
		q.Orchestration.Namespace,
		q.Orchestration.Name,
		time.Now().UTC(),
	).Scan(
		&oreID,
		&nsID,
		&createdAt,
		&createdBy,
		&since,
		&until,
		&namedAt,
		&namedBy,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	switch q.Orchestration.Type {
	case `server`:
		rows, err = h.stmtResolvNext.Query(
			oreID,
		)
	case `physical`:
		rows, err = h.stmtResolvPhys.Query(
			oreID,
		)
	default:
		mr.BadRequest()
		return
	}
	if err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&srvName,
			&nsName,
			&srvType,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		mr.ServerHeader = append(mr.ServerHeader, proto.ServerHeader{
			Namespace: nsName,
			Name:      srvName,
			Type:      srvType,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix