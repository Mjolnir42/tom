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
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdOrchestrationAdd)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationAdd,
		handle: orchestrationAdd,
	})
}

func orchestrationAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationAdd)
}

func exportOrchestrationAdd(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationAdd function
func (m *Model) OrchestrationAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = proto.ActionAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Orchestration = *req.Orchestration

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportOrchestrationAdd)
}

// add creates a new orchestration environment
func (h *OrchestrationWriteHandler) add(q *msg.Request, mr *msg.Result) {
	var (
		res                            sql.Result
		err                            error
		tx                             *sql.Tx
		txTime, validSince, validUntil time.Time
		rows                           *sql.Rows
		ok                             bool
		orchID                         string
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	if err = msg.ResolveValidSince(
		q.Orchestration.Property[`name`].ValidSince,
		&validSince, &txTime,
	); err != nil {
		mr.BadRequest(err)
		return
	}

	if err = msg.ResolveValidUntil(
		q.Orchestration.Property[`name`].ValidUntil,
		&validUntil, &txTime,
	); err != nil {
		mr.BadRequest(err)
		return
	}

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover all attributes and record them with their type
	attrMap := map[string]string{}
	if rows, err = tx.Query(
		stmt.NamespaceAttributeDiscover,
		q.Orchestration.Namespace,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	for rows.Next() {
		var attribute, typ string
		if err = rows.Scan(
			&attribute,
			&typ,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		attrMap[attribute] = typ
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// for all properties specified in the request, check that the attribute
	// exists and create missing attributes
	for key := range q.Orchestration.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = tx.Exec(
				stmt.NamespaceAttributeAddStandard,
				q.Orchestration.Namespace,
				key,
				q.UserIDLib,
				q.AuthUser,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
			if !mr.CheckRowsAffected(res.RowsAffected()) {
				tx.Rollback()
				return
			}
		}
	}

	// create named orchestration environment in specified namespace,
	// this is an INSERT statement with a RETURNING clause, thus
	// requires .QueryRow instead of .Exec
	if err = tx.Stmt(
		h.stmtAdd,
	).QueryRow(
		q.Orchestration.Namespace,
		q.UserIDLib,
		q.AuthUser,
		q.Orchestration.Property[`name`].Value,
		validSince,
		validUntil,
	).Scan(
		&orchID,
	); err == sql.ErrNoRows {
		// query did not return the generated rteID
		mr.ServerError(err)
		tx.Rollback()
		return
	} else if err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// for all properties specified in the request, update the value.
	// this transparently creates missing entries.
	for key := range q.Orchestration.Property {
		switch key {
		case `name`:
			continue
		default:
			if ok = h.txPropUpdate(
				q, mr, tx, &txTime, q.Orchestration.Property[key], orchID,
			); !ok {
				tx.Rollback()
				return
			}
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Orchestration = append(mr.Orchestration, q.Orchestration)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
