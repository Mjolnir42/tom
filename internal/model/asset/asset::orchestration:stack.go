/*-
 * Copyright (c) 2022, Jörg Pernfuß
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
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdOrchestrationStack)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationStack,
		handle: orchestrationStack,
	})
}

func orchestrationStack(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationStack)
}

func exportOrchestrationStack(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationStack function
func (m *Model) OrchestrationStack(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = proto.ActionStack

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Orchestration = *req.Orchestration

	request.Orchestration.TomID = params.ByName(`tomID`)
	if err := request.Orchestration.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	// input validation
	if err := proto.ValidNamespace(request.Orchestration.Namespace); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if err := proto.OnlyUnreserved(request.Orchestration.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if request.Orchestration.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}
	if len(request.Orchestration.Property) == 0 {
		m.x.ReplyBadRequest(&w, &request, nil)
	}

	for prop, obj := range request.Orchestration.Property {
		// property attribute must be stacking or unstacking
		// request
		switch request.Orchestration.Property[prop].Attribute {
		case proto.MetaPropertyCmdStack:
		case proto.MetaPropertyCmdUnstack:
		default:
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}

		// orchestrations must be provided-by: runtime
		if err, entityType, _ := proto.ParseTomID(
			request.Orchestration.Property[prop].Value,
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
	m.x.Send(&w, &result, exportOrchestrationStack)
}

// stack ...
func (h *OrchestrationWriteHandler) stack(q *msg.Request, mr *msg.Result) {
	var (
		oreID, nsID, rteID, rteNsID, createdAt, createdBy string
		nameValidSince, nameValidUntil, namedAt, namedBy  string
		err                                               error
		tx                                                *sql.Tx
	)

	// setup a consistent transaction time timestamp that is used for
	// all records
	txTime := time.Now().UTC()

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}
	defer tx.Rollback()

	txAdd := tx.Stmt(h.stmtTxStackAdd)
	txClamp := tx.Stmt(h.stmtTxStackClamp)
	txOreShow := tx.Stmt(h.stmtTxShow)
	txRteShow := tx.Stmt(h.stmtTxRteShow)

	// discover oreID at the start of the transaction
	if err = txOreShow.QueryRow(
		q.Orchestration.Namespace,
		q.Orchestration.Name,
		txTime,
	).Scan(
		&oreID,
		&nsID,
		&createdAt,
		&createdBy,
		&nameValidSince,
		&nameValidUntil,
		&namedAt,
		&namedBy,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	// handle all keys in the request
	for key := range q.Orchestration.Property {
		prop := q.Orchestration.Property[key]

		var reqValidSince, reqValidUntil time.Time
		if err = msg.ResolveValidSince(
			prop.ValidSince, &reqValidSince, &txTime,
		); err != nil {
			mr.ServerError(err)
			return
		}

		if err = msg.ResolveValidUntil(
			prop.ValidUntil, &reqValidUntil, &txTime,
		); err != nil {
			mr.ServerError(err)
			return
		}

		var ntt proto.Entity
		if err, _, ntt = proto.ParseTomID(
			prop.Value,
		); err != nil {
			mr.ServerError(err)
			return
		}

		if err = txRteShow.QueryRow(
			ntt.ExportNamespace(),
			ntt.ExportName(),
			reqValidSince,
		).Scan(
			&rteID,
			&rteNsID,
			&createdAt,
			&createdBy,
			&nameValidSince,
			&nameValidUntil,
			&namedAt,
			&namedBy,
		); err == sql.ErrNoRows {
			mr.NotFound(err)
			return
		} else if err != nil {
			mr.ServerError(err)
			return
		}

		switch prop.Attribute {
		case proto.MetaPropertyCmdStack:
			if err = h.txOreStack(
				txAdd,
				oreID, rteID,
				&txTime, &reqValidSince, &reqValidUntil,
				q, mr,
			); err != nil {
				mr.ServerError(err)
				return
			}
		case proto.MetaPropertyCmdUnstack:
			if err = h.txOreClamp(
				txClamp,
				oreID, rteID,
				&txTime, &reqValidSince, &reqValidUntil,
				q, mr,
			); err != nil {
				mr.ServerError(err)
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

func (h *OrchestrationWriteHandler) txOreStack(
	txStmt *sql.Stmt,
	oreID, rteID string,
	txTime, since, until *time.Time,
	q *msg.Request,
	mr *msg.Result,
) error {
	var err error
	var res sql.Result

	// create new record
	if res, err = txStmt.Exec(
		oreID,  // orchestration
		rteID,  // parent runtime
		since,  // lower(validity)
		until,  // upper(validity)
		txTime, // createdAt
		q.AuthUser,
		q.UserIDLib,
	); err != nil {
		mr.ServerError(err)
		return err
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return mr.ExportError()
	}
	return err
}

func (h *OrchestrationWriteHandler) txOreClamp(
	txStmt *sql.Stmt,
	oreID, rteID string,
	txTime, since, until *time.Time,
	q *msg.Request,
	mr *msg.Result,
) error {
	var err error
	var res sql.Result

	// clamp existing records
	if res, err = txStmt.Exec(
		since, // this is used as `invalid since` / `valid until`
		oreID,
		rteID,  // parent runtime to clamp
		txTime, // parent relationship must be valid at txTime
	); err != nil {
		mr.ServerError(err)
		return err
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return mr.ExportError()
	}
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix