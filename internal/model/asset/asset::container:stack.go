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
	proto.AssertCommandIsDefined(proto.CmdContainerStack)

	registry = append(registry, function{
		cmd:    proto.CmdContainerStack,
		handle: containerStack,
	})
}

func containerStack(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ContainerStack)
}

func exportContainerStack(result *proto.Result, r *msg.Result) {
	result.Container = &[]proto.Container{}
	*result.Container = append(*result.Container, r.Container...)
}

// ContainerStack function
func (m *Model) ContainerStack(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdContainerStack,
		msg.SectionContainer,
		proto.ActionStack,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Container = *req.Container

	request.Container.TomID = params.ByName(`tomID`)
	if err := request.Container.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	// input validation
	if err := proto.ValidNamespace(request.Container.Namespace); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if err := proto.OnlyUnreserved(request.Container.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if request.Container.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}
	if len(request.Container.Property) != 1 {
		m.x.ReplyBadRequest(&w, &request, nil)
	}

	for prop, obj := range request.Container.Property {
		// property attribute must be stacking request
		switch request.Container.Property[prop].Attribute {
		case proto.MetaPropertyCmdStack:
		default:
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}

		// container must be provided-by: runtime
		if err, entityType, _ := proto.ParseTomID(
			request.Container.Property[prop].Value,
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
	m.x.Send(&w, &result, exportContainerStack)
}

// stack ...
func (h *ContainerWriteHandler) stack(q *msg.Request, mr *msg.Result) {
	var (
		cnID, nsID, rteID, rteNsID, createdAt, createdBy string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		err                                              error
		tx                                               *sql.Tx
		res                                              sql.Result
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
	txShow := tx.Stmt(h.stmtTxShow)
	txRteShow := tx.Stmt(h.stmtTxRteShow)

	// discover cnID at the start of the transaction
	if err = txShow.QueryRow(
		q.Container.Namespace,
		q.Container.Name,
		txTime,
	).Scan(
		&cnID,
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

	// there is exactly one entry, but we do not know the key
	for key := range q.Container.Property {
		prop := q.Container.Property[key]

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

		// clamp validity of existing records to the lower bound
		// of the new record
		if _, err = txClamp.Exec(
			reqValidSince,
			cnID,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}

		// create new record
		if res, err = txAdd.Exec(
			cnID,
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
	mr.Container = append(mr.Container, q.Container)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
