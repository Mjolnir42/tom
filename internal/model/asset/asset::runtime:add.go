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
	proto.AssertCommandIsDefined(proto.CmdRuntimeAdd)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeAdd,
		handle: runtimeAdd,
	})
}

func runtimeAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeAdd)
}

func exportRuntimeAdd(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimeAdd function
func (m *Model) RuntimeAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionRuntime
	request.Action = proto.ActionAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Runtime = *req.Runtime

	if err := proto.ValidNamespace(
		request.Runtime.Namespace,
	); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	// at least the name property must be filled in
	if request.Runtime.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}

	for prop, obj := range request.Runtime.Property {
		if err := proto.OnlyUnreserved(prop); err != nil {
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
	m.x.Send(&w, &result, exportRuntimeAdd)
}

// add creates a new runtime environment
func (h *RuntimeWriteHandler) add(q *msg.Request, mr *msg.Result) {
	var (
		res                            sql.Result
		err                            error
		tx                             *sql.Tx
		txTime, validSince, validUntil time.Time
		rows                           *sql.Rows
		ok                             bool
		rteID                          string
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	switch q.Runtime.Property[`name`].ValidSince {
	case `always`:
		validSince = msg.NegTimeInf
	case `forever`:
		mr.BadRequest()
		return
	case ``:
		validSince = txTime
	default:
		if validSince, err = time.Parse(
			msg.RFC3339Milli, q.Runtime.Property[`name`].ValidSince,
		); err != nil {
			mr.BadRequest(err)
			return
		}
	}

	switch q.Runtime.Property[`name`].ValidUntil {
	case `always`:
		mr.BadRequest()
		return
	case `forever`:
		validUntil = msg.PosTimeInf
	case ``:
		validUntil = msg.PosTimeInf
	default:
		if validSince, err = time.Parse(
			msg.RFC3339Milli, q.Runtime.Property[`name`].ValidUntil,
		); err != nil {
			mr.BadRequest(err)
			return
		}
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
		q.Runtime.Namespace,
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
	for key := range q.Runtime.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = tx.Exec(
				stmt.NamespaceAttributeAddStandard,
				q.Runtime.Namespace,
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

	// create named runtime environment in specified namespace,
	// this is an INSERT statement with a RETURNING clause, thus
	// requires .QueryRow instead of .Exec
	if err = tx.Stmt(h.stmtAdd).QueryRow(
		q.Runtime.Namespace,
		q.UserIDLib,
		q.AuthUser,
		q.Runtime.Property[`name`].Value,
		validSince,
		validUntil,
	).Scan(
		&rteID,
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
	for key := range q.Runtime.Property {
		if key == `name` {
			continue
		}
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Runtime.Property[key], rteID,
		); !ok {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Runtime = append(mr.Runtime, q.Runtime)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
