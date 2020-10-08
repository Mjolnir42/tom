/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package server // import "github.com/mjolnir42/tom/internal/core/asset/server"

import (
	"database/sql"

	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

// process is the request dispatcher
func (h *ReadHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		h.list(q, &result)
	case msg.ActionShow:
		h.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all servers
func (h *ReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		id, namespace, name, typ string
		rows                     *sql.Rows
		err                      error
	)

	if rows, err = h.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&namespace,
			&name,
			&typ,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		mr.Server = append(mr.Server, proto.Server{
			ID:        id,
			Namespace: namespace,
			Name:      name,
			Type:      typ,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// show returns full details for a specific server
func (h *ReadHandler) show(q *msg.Request, mr *msg.Result) {
	/* TODO
	var (
		id, namespace, name, typ string
		err                      error
	)
	if err = h.stmtShow.QueryRow(
		q.Server.ID,
	).Scan(
		&id,
		&namespace,
		&name,
		&typ,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Server = append(mr.Server, proto.Server{
		ID:        id,
		Namespace: namespace,
		Name:      name,
		Type:      typ,
	})
	*/
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
