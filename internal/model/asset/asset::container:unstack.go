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
	proto.AssertCommandIsDefined(proto.CmdContainerUnstack)

	registry = append(registry, function{
		cmd:    proto.CmdContainerUnstack,
		handle: containerUnstack,
	})
}

func containerUnstack(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ContainerUnstack)
}

func exportContainerUnstack(result *proto.Result, r *msg.Result) {
	result.Container = &[]proto.Container{}
	*result.Container = append(*result.Container, r.Container...)
}

// ContainerUnstack function
func (m *Model) ContainerUnstack(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdContainerUnstack,
		msg.SectionContainer,
		proto.ActionUnstack,
	)

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

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportContainerUnstack)
}

// unstack ...
func (h *ContainerWriteHandler) unstack(q *msg.Request, mr *msg.Result) {
	var (
		cnID, nsID, createdAt, createdBy                 string
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

	// discover oreID at the start of the transaction
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

	if _, err = txClamp.Exec(
		txTime,
		cnID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Container = append(mr.Container, q.Container)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
