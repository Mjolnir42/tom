/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
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
	proto.AssertCommandIsDefined(proto.CmdOrchestrationPropUpdate)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationPropUpdate,
		handle: orchestrationPropUpdate,
	})
}

func orchestrationPropUpdate(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationPropUpdate)
}

func exportOrchestrationPropUpdate(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationPropUpdate function
func (m *Model) OrchestrationPropUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdOrchestrationPropUpdate,
		msg.SectionOrchestration,
		proto.ActionPropUpdate,
	)

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
	for prop, obj := range request.Orchestration.Property {
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
	m.x.Send(&w, &result, exportOrchestrationPropUpdate)
}

// propertyUpdate ...
func (h *OrchestrationWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	var (
		oreID, nsID, createdAt, createdBy                string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		err                                              error
		tx                                               *sql.Tx
		rows                                             *sql.Rows
		res                                              sql.Result
		ok                                               bool
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

	txAttAddStd := tx.Stmt(h.stmtAttAddStandard)
	txDiscover := tx.Stmt(h.stmtAttDiscover)
	txShow := tx.Stmt(h.stmtTxShow)

	// discover oreID at the start of the transaction, as the property
	// updates might include a name change
	if err = txShow.QueryRow(
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

	// discover all attributes available within the namespace
	// and record them with their type
	attrMap := map[string]string{}
	if rows, err = txDiscover.Query(
		q.Orchestration.Namespace,
	); err != nil {
		mr.ServerError(err)
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
			return
		}
		attrMap[attribute] = typ
	}

	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// for all properties specified in the request, check that the attribute
	// exists and create missing attributes
	for key := range q.Orchestration.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = txAttAddStd.Exec(
				q.Orchestration.Namespace,
				key,
				q.UserIDLib,
				q.AuthUser,
			); err != nil {
				mr.ServerError(err)
				return
			}
			if !mr.CheckRowsAffected(res.RowsAffected()) {
				tx.Rollback()
				return
			}
		}
	}

	// for all properties specified in the request, update the value.
	// this transparently creates missing values.
	for key := range q.Orchestration.Property {
		if key == `type` || q.Orchestration.Property[key].Attribute == `type` {
			continue
		}
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Orchestration.Property[key], oreID,
		); !ok {
			tx.Rollback()
			return
		}
		// remove the property from the map of available attributes
		delete(attrMap, q.Orchestration.Property[key].Attribute)
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Orchestration = append(mr.Orchestration, q.Orchestration)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
