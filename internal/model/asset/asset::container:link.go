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
	proto.AssertCommandIsDefined(proto.CmdContainerLink)

	registry = append(registry, function{
		cmd:    proto.CmdContainerLink,
		handle: containerLink,
	})
}

func containerLink(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ContainerLink)
}

func exportContainerLink(result *proto.Result, r *msg.Result) {
	result.Container = &[]proto.Container{}
	*result.Container = append(*result.Container, r.Container...)
}

// ContainerPropUpdate function
func (m *Model) ContainerLink(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdContainerLink,
		msg.SectionContainer,
		proto.ActionLink,
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
	for prop, obj := range request.Container.Property {
		if request.Container.Property[prop].Attribute != proto.MetaPropertyCmdLink {
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}
		if err, entityType, _ := proto.ParseTomID(
			request.Container.Property[prop].Value,
		); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		} else if entityType != proto.EntityContainer {
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
	m.x.Send(&w, &result, exportContainerLink)
}

// link ...
func (h *ContainerWriteHandler) link(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                            time.Time
		containerID, dictionaryID, linkContID, linkDictID string
		createdAt, createdBy                              string
		nameValidSince, nameValidUntil, namedAt, namedBy  string
		err                                               error
		tx                                                *sql.Tx
		res                                               sql.Result
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover containerID at the start of the transaction, as the property
	// updates might include a name change
	if err = tx.QueryRow(
		stmt.ContainerTxShow,
		q.Container.Namespace,
		q.Container.Name,
		txTime,
	).Scan(
		&containerID,
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
	for key := range q.Container.Property {
		ct := proto.Container{
			TomID: q.Container.Property[key].Value,
		}
		if err = ct.ParseTomID(); err != nil {
			mr.BadRequest(err)
			tx.Rollback()
			return
		}
		// query IDs for link target
		if err = tx.QueryRow(
			stmt.ContainerTxShow,
			ct.Namespace,
			ct.Name,
			txTime,
		).Scan(
			&linkContID,
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
		var contUUID, linkUUID uuid.UUID
		if contUUID, err = uuid.FromString(containerID); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		if linkUUID, err = uuid.FromString(linkContID); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		switch bytes.Compare(contUUID.Bytes(), linkUUID.Bytes()) {
		case 0: // contUUID == linkUUID
			mr.ServerError(fmt.Errorf("Cannot link container environments to itself."))
			tx.Rollback()
			return
		case -1: // contUUID < linkUUID
			if res, err = tx.Stmt(h.stmtLink).Exec(
				linkContID,
				linkDictID,
				containerID,
				dictionaryID,
				q.AuthUser,
				q.UserIDLib,
				txTime,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
		case 1: // linkUUID < contUUID
			if res, err = tx.Stmt(h.stmtLink).Exec(
				containerID,
				dictionaryID,
				linkContID,
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
			mr.ServerError(fmt.Errorf("Error ordering container IDs for linking"))
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
	mr.Container = append(mr.Container, q.Container)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
