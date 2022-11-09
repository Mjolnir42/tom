/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdServerStack)

	registry = append(registry, function{
		cmd:    proto.CmdServerStack,
		handle: serverStack,
	})
}

func serverStack(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ServerStack)
}

func exportServerStack(result *proto.Result, r *msg.Result) {
	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)
}

// ServerStack function
func (m *Model) ServerStack(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdServerStack,
		msg.SectionServer,
		proto.ActionStack,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Server = *req.Server

	request.Server.TomID = params.ByName(`tomID`)
	if err := request.Server.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	// input validation
	if err := proto.ValidNamespace(request.Server.Namespace); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if err := proto.OnlyUnreserved(request.Server.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if request.Server.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}
	if len(request.Server.Property) != 1 {
		m.x.ReplyBadRequest(&w, &request, nil)
	}

	for prop, obj := range request.Server.Property {
		// property attribute must be stacking request
		switch request.Server.Property[prop].Attribute {
		case proto.MetaPropertyCmdStack:
		default:
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}

		// servers must be provided-by runtimes
		if err, entityType, _ := proto.ParseTomID(
			request.Server.Property[prop].Value,
		); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		} else {
			switch entityType {
			case proto.EntityRuntime:
			default:
				m.x.ReplyBadRequest(&w, &request, nil)
				return
			}
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
	m.x.Send(&w, &result, exportServerStack)
}

// stack ...
func (h *ServerWriteHandler) stack(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		serverID, dictionaryID, createdAt, createdBy     string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		rows                                             *sql.Rows
		err                                              error
		tx                                               *sql.Tx
		res                                              sql.Result
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover serverID at the start of the transaction, as the property
	// updates might include a name change
	if err = tx.QueryRow(
		stmt.ServerTxShow,
		q.Server.Namespace,
		q.Server.Name,
		txTime,
	).Scan(
		&serverID,
		&dictionaryID,
		&createdAt,
		&createdBy,
		&nameValidSince,
		&nameValidUntil,
		&namedAt,
		&namedBy,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		tx.Rollback()
		return
	} else if err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// check server type property - only virtual may be stacked
	txProp := tx.Stmt(h.stmtTxProp)

	if rows, err = txProp.Query(
		dictionaryID,
		serverID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	for rows.Next() {
		prop := proto.PropertyDetail{}

		if err = rows.Scan(
			&prop.Attribute,
			&prop.Value,
			&prop.ValidSince,
			&prop.ValidUntil,
			&prop.CreatedAt,
			&prop.CreatedBy,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		switch prop.Attribute {
		case `type`:
			if prop.Value != `virtual` {
				rows.Close()
				mr.BadRequest(fmt.Errorf(
					`Illegal request to stack a non-virtual server`,
				))
				tx.Rollback()
				return
			}
		default:
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	txAdd := tx.Stmt(h.stmtTxStackAdd)
	txClamp := tx.Stmt(h.stmtTxStackClamp)
	txRte := tx.Stmt(h.stmtTxRuntimeShow)

	// only one key, but we can't be sure of its value
	for key := range q.Server.Property {
		prop := q.Server.Property[key]

		// check the use of the perpetual keyword
		if prop.ValidSince == `perpetual` || prop.ValidUntil == `perpetual` {
			// only the `type` property is allowed to be perpetual
			mr.BadRequest()
			tx.Rollback()
			return
		}

		var reqValidSince, reqValidUntil time.Time
		var ntt proto.Entity
		var rteID string

		switch prop.ValidSince {
		case `always`:
			reqValidSince = msg.NegTimeInf
		case `forever`:
			mr.BadRequest()
			tx.Rollback()
			return
		case ``, `now`:
			reqValidSince = txTime
		default:
			if reqValidSince, err = time.Parse(
				msg.RFC3339Milli,
				prop.ValidSince,
			); err != nil {
				mr.BadRequest(err)
				tx.Rollback()
				return
			}
		}

		switch prop.ValidUntil {
		case `always`:
			mr.BadRequest()
			tx.Rollback()
			return
		case `forever`:
			reqValidUntil = msg.PosTimeInf
		case ``:
			reqValidUntil = msg.PosTimeInf
		case `now`:
			reqValidUntil = txTime
		default:
			if reqValidUntil, err = time.Parse(
				msg.RFC3339Milli,
				prop.ValidUntil,
			); err != nil {
				mr.BadRequest(err)
				tx.Rollback()
				return
			}
		}

		// clamp validity of existing records to the lower bound of the new
		// record
		if _, err = txClamp.Exec(
			reqValidSince,
			serverID,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		// resolve rteID
		if err, _, ntt = proto.ParseTomID(
			prop.Value,
		); err == sql.ErrNoRows {
			mr.NotFound(err)
			tx.Rollback()
			return
		} else if err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		if err = txRte.QueryRow(
			ntt.ExportNamespace(),
			ntt.ExportName(),
			reqValidSince,
		).Scan(
			&rteID,
			&dictionaryID,
			&createdAt,
			&createdBy,
			&nameValidSince,
			&nameValidUntil,
			&namedAt,
			&namedBy,
		); err == sql.ErrNoRows {
			mr.NotFound(err)
			tx.Rollback()
			return
		} else if err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		// create new record
		if res, err = txAdd.Exec(
			serverID,
			rteID,
			reqValidSince, // lower(validity)
			reqValidUntil, // upper(validity)
			txTime,        // createdAt
			q.AuthUser,
			q.UserIDLib,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		if !mr.AssertOneRowAffected(res.RowsAffected()) {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Server = append(mr.Server, q.Server)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
