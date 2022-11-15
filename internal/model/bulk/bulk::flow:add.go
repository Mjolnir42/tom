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
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdFlowAdd)

	registry = append(registry, function{
		cmd:    proto.CmdFlowAdd,
		handle: flowAdd,
	})
}

func flowAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.FlowAdd)
}

func exportFlowAdd(result *proto.Result, r *msg.Result) {
	result.Flow = &[]proto.Flow{}
	*result.Flow = append(*result.Flow, r.Flow...)
}

// FlowAdd function
func (m *Model) FlowAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(
		r, params,
		proto.CmdFlowAdd,
		msg.SectionFlow,
		proto.ActionAdd,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Flow = *req.Flow

	if err := proto.ValidNamespace(
		request.Flow.Namespace,
	); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	// at least the name property must be filled in
	if request.Flow.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}

	for prop, obj := range request.Flow.Property {
		if err := proto.ValidAttribute(prop); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
		if err := proto.CheckPropertyConstraints(&obj); err != nil {
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
	m.x.Send(&w, &result, exportFlowAdd)
}

// add creates a new flow
func (h *FlowWriteHandler) add(q *msg.Request, mr *msg.Result) {
	var (
		res                            sql.Result
		err                            error
		tx                             *sql.Tx
		txTime, validSince, validUntil time.Time
		rows                           *sql.Rows
		ok                             bool
		flowID                         string
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	if err = msg.ResolveValidSince(
		q.Flow.Property[`name`].ValidSince,
		&validSince, &txTime,
	); err != nil {
		mr.BadRequest(err)
		return
	}

	if err = msg.ResolveValidUntil(
		q.Flow.Property[`name`].ValidUntil,
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
		q.Flow.Namespace,
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
	for key := range q.Flow.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = tx.Exec(
				stmt.NamespaceAttributeAddStandard,
				q.Flow.Namespace,
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

	// create named flow in specified namespace,
	// this is an INSERT statement with a RETURNING clause, thus
	// requires .QueryRow instead of .Exec
	if err = tx.Stmt(h.stmtAdd).QueryRow(
		q.Flow.Namespace,
		q.UserIDLib,
		q.AuthUser,
		q.Flow.Property[`name`].Value,
		validSince,
		validUntil,
	).Scan(
		&flowID,
	); err == sql.ErrNoRows {
		// query did not return the generated flowID
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
	for key := range q.Flow.Property {
		if key == `name` {
			continue
		}
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Flow.Property[key], flowID,
		); !ok {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Flow = append(mr.Flow, q.Flow)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
