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
	proto.AssertCommandIsDefined(proto.CmdServerLink)

	registry = append(registry, function{
		cmd:    proto.CmdServerLink,
		handle: serverLink,
	})
}

func serverLink(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ServerLink)
}

func exportServerLink(result *proto.Result, r *msg.Result) {
	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)
}

// ServerLink function
func (m *Model) ServerLink(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = proto.ActionLink

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
	for prop, obj := range request.Server.Property {
		if request.Server.Property[prop].Attribute != proto.MetaPropertyCmdLink {
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}
		if err, entityType, _ := proto.ParseTomID(
			request.Server.Property[prop].Value,
		); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		} else if entityType != proto.EntityServer {
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
	m.x.Send(&w, &result, exportServerLink)
}

// link ...
func (h *ServerWriteHandler) link(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		serverID, dictionaryID, linkServerID, linkDictID string
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

	// for all properties specified in the request, update the value.
	// this transparently creates missing entries.
	for key := range q.Server.Property {
		ct := proto.Server{
			TomID: q.Server.Property[key].Value,
		}
		if err = ct.ParseTomID(); err != nil {
			mr.BadRequest(err)
			tx.Rollback()
			return
		}
		// query IDs for link target
		if err = tx.QueryRow(
			stmt.ServerTxShow,
			ct.Namespace,
			ct.Name,
			txTime,
		).Scan(
			&linkServerID,
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
		var srvUUID, linkUUID uuid.UUID
		if srvUUID, err = uuid.FromString(serverID); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		if linkUUID, err = uuid.FromString(linkServerID); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		switch bytes.Compare(srvUUID.Bytes(), linkUUID.Bytes()) {
		case 0: // srvUUID == linkUUID
			mr.ServerError(fmt.Errorf("Cannot link server to itself."))
			tx.Rollback()
			return
		case -1: // srvUUID < linkUUID
			if res, err = tx.Stmt(h.stmtLink).Exec(
				linkServerID,
				linkDictID,
				serverID,
				dictionaryID,
				q.AuthUser,
				q.UserIDLib,
				txTime,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
		case 1: // linkUUID < srvUUID
			if res, err = tx.Stmt(h.stmtLink).Exec(
				serverID,
				dictionaryID,
				linkServerID,
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
			mr.ServerError(fmt.Errorf("Error ordering server IDs for linking"))
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
	mr.Server = append(mr.Server, q.Server)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
