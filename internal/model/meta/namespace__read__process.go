/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"time"

	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

// process is the request dispatcher
func (h *NamespaceReadHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

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

// list returns all namespaces
func (h *NamespaceReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		dictionaryName, author string
		creationTime           time.Time
		rows                   *sql.Rows
		err                    error
	)

	if rows, err = h.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&dictionaryName,
			&creationTime,
			&author,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		mr.NamespaceHeader = append(mr.NamespaceHeader, proto.NamespaceHeader{
			Name:      dictionaryName,
			CreatedAt: creationTime.Format(msg.RFC3339Milli),
			CreatedBy: author,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// show returns full details for a specific server
func (h *NamespaceReadHandler) show(q *msg.Request, mr *msg.Result) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
