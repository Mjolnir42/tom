/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/satori/go.uuid"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdRuntimeLink)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeLink,
		handle: runtimeLink,
	})
}

func runtimeLink(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeLink)
}

func exportRuntimeLink(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimePropUpdate function
func (m *Model) RuntimeLink(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionRuntime
	request.Action = proto.ActionLink

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
	for prop, obj := range request.Runtime.Property {
		if request.Runtime.Property[prop].Attribute != proto.ActionLink {
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}
		if err := proto.OnlyUnreserved(
			request.Runtime.Property[prop].Attribute,
		); err != nil {
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
	m.x.Send(&w, &result, exportRuntimeLink)
}

// link ...
func (h *RuntimeWriteHandler) link(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		rteID, dictionaryID, linkRteID, linkDictID       string
		createdAt, createdBy                             string
		nameValidSince, nameValidUntil, namedAt, namedBy string
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

	// discover rteID at the start of the transaction, as the property
	// updates might include a name change
	if err = tx.QueryRow(
		stmt.RuntimeTxShow,
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

	// for all properties specified in the request, update the value.
	// this transparently creates missing entries.
	for key := range q.Runtime.Property {
		rte := proto.Runtime{
			TomID: q.Runtime.Property[key].Value,
		}
		if err = rte.ParseTomID(); err != nil {
			mr.BadRequest(err)
			tx.Rollback()
			return
		}
		// query IDs for link target
		if err = tx.QueryRow(
			stmt.RuntimeTxShow,
			rte.Namespace,
			rte.Name,
			txTime,
		).Scan(
			&linkRteID,
			&linkDictID,
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
		// convert UUID strings to UUID objects to access their byte arrays
		var rteUUID, linkUUID uuid.UUID
		if rteUUID, err = uuid.FromString(rteID); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		if linkUUID, err = uuid.FromString(linkRteID); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		switch bytes.Compare(rteUUID.Bytes(), linkUUID.Bytes()) {
		case 0: // rteUUID == linkUUID
			mr.ServerError(fmt.Errorf("Cannot link runtime environments to itself."))
			tx.Rollback()
			return
		case -1: // rteUUID < linkUUID
			if res, err = tx.Stmt(h.stmtLink).Exec(
				linkRteID,
				linkDictID,
				rteID,
				dictionaryID,
				q.AuthUser,
				q.UserIDLib,
				txTime,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
		case 1: // linkUUID < rteUUID
			if res, err = tx.Stmt(h.stmtLink).Exec(
				rteID,
				dictionaryID,
				linkRteID,
				linkDictID,
				q.AuthUser,
				q.UserIDLib,
				txTime,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
		default:
			mr.ServerError(fmt.Errorf("Error ordering runtime IDs for linking"))
			tx.Rollback()
			return
		}
		if !mr.CheckRowsAffected(res.RowsAffected()) {
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
