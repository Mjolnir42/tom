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
	proto.AssertCommandIsDefined(proto.CmdRuntimeResolve)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeResolve,
		handle: runtimeResolve,
	})
}

func runtimeResolve(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeResolve)
}

func exportRuntimeResolve(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	result.ServerHeader = &[]proto.ServerHeader{}
	*result.ServerHeader = append(*result.ServerHeader, r.ServerHeader...)
}

// RuntimeResolve function
func (m *Model) RuntimeResolve(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionRuntime
	request.Action = proto.ActionResolve
	request.Runtime = proto.Runtime{
		TomID: params.ByName(`tomID`),
		Type:  params.ByName(`level`), // resolution detail type
	}

	if err := request.Runtime.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportRuntimeResolve)
}

// resolve ...
func (h *RuntimeReadHandler) resolve(q *msg.Request, mr *msg.Result) {
	var (
		nsName, nsID, rteID, srvName, srvType  string
		createdAt, createdBy, namedAt, namedBy string
		since, until                           time.Time
		rows                                   *sql.Rows
		err                                    error
	)

	if err = h.stmtShow.QueryRow(
		q.Runtime.Namespace,
		q.Runtime.Name,
		time.Now().UTC(),
	).Scan(
		&rteID,
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

	switch q.Runtime.Type {
	case `server`:
		rows, err = h.stmtResolvNext.Query(
			rteID,
		)
	case `physical`:
		rows, err = h.stmtResolvPhys.Query(
			rteID,
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
