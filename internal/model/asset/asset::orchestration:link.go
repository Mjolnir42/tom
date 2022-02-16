/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
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
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/satori/go.uuid"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdOrchestrationLink)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationLink,
		handle: orchestrationLink,
	})
}

func orchestrationLink(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationLink)
}

func exportOrchestrationLink(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationLink  function
func (m *Model) OrchestrationLink(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdOrchestrationLink,
		msg.SectionOrchestration,
		proto.ActionLink,
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
		if request.Orchestration.Property[prop].Attribute != proto.MetaPropertyCmdLink {
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}
		if err, entityType, _ := proto.ParseTomID(
			request.Orchestration.Property[prop].Value,
		); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		} else if entityType != proto.EntityOrchestration {
			m.x.ReplyBadRequest(&w, &request, nil)
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
	m.x.Send(&w, &result, exportOrchestrationLink)
}

// link ...
func (h *OrchestrationWriteHandler) link(q *msg.Request, mr *msg.Result) {
	var (
		oreID, nsID, createdAt, createdBy                string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		tx                                               *sql.Tx
		err                                              error
		res                                              sql.Result
	)

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}
	defer tx.Rollback()

	// setup a consistent transaction time timestamp that is used
	// for all records
	txTime := time.Now().UTC()

	txShow := tx.Stmt(h.stmtTxShow)
	txLink := tx.Stmt(h.stmtTxLink)

	// discover oreID
	if err = txShow.QueryRow(
		q.Orchestration.Namespace,
		q.Orchestration.Name,
		txTime,
	).Scan(
		&oreID,
		&nsID,
		&createdAt,      // unused result
		&createdBy,      // unused result
		&nameValidSince, // unused result
		&nameValidUntil, // unused result
		&namedAt,        // unused result
		&namedBy,        // unused result
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	// for targets specified in the request, update the links
	for key := range q.Orchestration.Property {
		var linkOreID, linkNsID string
		ore := proto.Orchestration{
			TomID: q.Orchestration.Property[key].Value,
		}
		if err = ore.ParseTomID(); err != nil {
			mr.BadRequest(err)
			return
		}
		// query IDs for link target
		if err = txShow.QueryRow(
			ore.Namespace,
			ore.Name,
			txTime,
		).Scan(
			&linkOreID,
			&linkNsID,
			&createdAt,      // unused
			&createdBy,      // unused
			&nameValidSince, // unused
			&nameValidUntil, // unused
			&namedAt,        // unused
			&namedBy,        // unused
		); err == sql.ErrNoRows {
			mr.NotFound(err)
			return
		} else if err != nil {
			mr.ServerError(err)
			return
		}

		// convert UUID strings to UUID objects to access their byte arrays
		var oreUUID, linkUUID uuid.UUID
		if oreUUID, err = uuid.FromString(oreID); err != nil {
			mr.ServerError(err)
			return
		}
		if linkUUID, err = uuid.FromString(linkOreID); err != nil {
			mr.ServerError(err)
			return
		}

		switch bytes.Compare(oreUUID.Bytes(), linkUUID.Bytes()) {
		case 0: // oreUUID == linkUUID
			mr.ServerError(
				fmt.Errorf("Cannot link orchestration with itself."))
			return
		case -1: // oreUUID < linkUUID
			if res, err = txLink.Exec(
				linkOreID,
				linkNsID,
				oreID,
				nsID,
				q.AuthUser,
				q.UserIDLib,
				txTime,
			); err != nil {
				mr.ServerError(err)
				return
			}
		case 1: // linkUUID < oreUUID
			if res, err = txLink.Exec(
				oreID,
				nsID,
				linkOreID,
				linkNsID,
				q.AuthUser,
				q.UserIDLib,
				txTime,
			); err != nil {
				mr.ServerError(err)
				return
			}
		default:
			mr.ServerError(
				fmt.Errorf("Error ordering orchestration IDs for linking"))
			return
		}
		if !mr.CheckRowsAffected(res.RowsAffected()) {
			return
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
