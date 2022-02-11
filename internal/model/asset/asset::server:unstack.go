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
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdServerUnstack)

	registry = append(registry, function{
		cmd:    proto.CmdServerUnstack,
		handle: serverUnstack,
	})
}

func serverUnstack(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ServerUnstack)
}

func exportServerUnstack(result *proto.Result, r *msg.Result) {
	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)
}

// ServerUnstack function
func (m *Model) ServerUnstack(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = proto.ActionUnstack
	request.Server = *(proto.NewServer())

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

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportServerUnstack)
}

// unstack ...
func (h *ServerWriteHandler) unstack(q *msg.Request, mr *msg.Result) {
	var (
		srvID, nsID, createdAt, createdBy                string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		err                                              error
		tx                                               *sql.Tx
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

	txClamp := tx.Stmt(h.stmtTxStackClamp)
	txShow := tx.Stmt(h.stmtTxShow)

	// discover srvID at the start of the transaction
	if err = txShow.QueryRow(
		q.Server.Namespace,
		q.Server.Name,
		txTime,
	).Scan(
		&srvID,
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

	if _, err = txClamp.Exec(
		txTime,
		srvID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Server = append(mr.Server, q.Server)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
