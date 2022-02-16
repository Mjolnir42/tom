/*-
 * Copyright (c) 2022, Jörg Pernfuß
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
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdRuntimeStack)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeStack,
		handle: runtimeStack,
	})
}

func runtimeStack(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeStack)
}

func exportRuntimeStack(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimeStack function
func (m *Model) RuntimeStack(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdRuntimeStack,
		msg.SectionRuntime,
		proto.ActionStack,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Runtime = *req.Runtime

	request.Runtime.TomID = params.ByName(`tomID`)
	if err := request.Runtime.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	// input validation
	if err := proto.ValidNamespace(request.Runtime.Namespace); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if err := proto.OnlyUnreserved(request.Runtime.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if request.Runtime.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}
	if len(request.Runtime.Property) != 1 {
		m.x.ReplyBadRequest(&w, &request, nil)
	}

	for prop, obj := range request.Runtime.Property {
		// property attribute must be stacking request
		if request.Runtime.Property[prop].Attribute != proto.MetaPropertyCmdStack {
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}
		// runtime must be runs-on: runtime, orchestration, server
		if err, entityType, _ := proto.ParseTomID(
			request.Runtime.Property[prop].Value,
		); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		} else {
			switch entityType {
			case proto.EntityRuntime:
			case proto.EntityOrchestration:
			case proto.EntityServer:
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
	m.x.Send(&w, &result, exportRuntimeStack)
}

// stack ...
func (h *RuntimeWriteHandler) stack(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		rteID, dictionaryID, createdAt, createdBy        string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		parentRteID, parentServerID, parentOrchID        sql.NullString
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

	txAdd := tx.Stmt(h.stmtTxStackAdd)
	txClamp := tx.Stmt(h.stmtTxStackClamp)
	txRteShow := tx.Stmt(h.stmtTxShow)
	txSrvShow := tx.Stmt(h.stmtTxServerShow)
	txOrcShow := tx.Stmt(h.stmtTxOrchShow)

	// discover rteID at the start of the transaction
	if err = txRteShow.QueryRow(
		q.Runtime.Namespace,
		q.Runtime.Name,
		txTime,
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

	// only one key, but we can't be sure of its value
	for key := range q.Runtime.Property {
		prop := q.Runtime.Property[key]

		// check the use of the perpetual keyword
		if prop.ValidSince == `perpetual` || prop.ValidUntil == `perpetual` {
			// only the `type` property is allowed to be perpetual
			mr.BadRequest()
			tx.Rollback()
			return
		}

		var reqValidSince, reqValidUntil time.Time
		var ntt proto.Entity
		var nttType string

		switch prop.ValidSince {
		case `always`:
			reqValidSince = msg.NegTimeInf
		case `forever`:
			mr.BadRequest()
			tx.Rollback()
			return
		case ``:
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
			rteID,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		// resolve new parent entity ID
		if err, nttType, ntt = proto.ParseTomID(
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

		switch nttType {
		case proto.EntityRuntime:
			if err = txRteShow.QueryRow(
				ntt.ExportNamespace(),
				ntt.ExportName(),
				reqValidSince,
			).Scan(
				&(parentRteID.String),
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
			parentRteID.Valid = true
		case proto.EntityServer:
			if err = txSrvShow.QueryRow(
				ntt.ExportNamespace(),
				ntt.ExportName(),
				reqValidSince,
			).Scan(
				&(parentServerID.String),
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
			parentServerID.Valid = true
		case proto.EntityOrchestration:
			if err = txOrcShow.QueryRow(
				ntt.ExportNamespace(),
				ntt.ExportName(),
				reqValidSince,
			).Scan(
				&(parentOrchID.String),
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
			parentOrchID.Valid = true
		default:
			mr.BadRequest(fmt.Errorf(
				"Invalid parent object class: %s", nttType,
			))
			tx.Rollback()
			return
		}

		// create new record
		if res, err = txAdd.Exec(
			rteID,
			parentServerID,
			parentRteID,
			parentOrchID,
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
	mr.Runtime = append(mr.Runtime, q.Runtime)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
