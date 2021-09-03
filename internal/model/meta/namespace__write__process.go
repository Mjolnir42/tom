/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"time"

	"github.com/mjolnir42/tom/internal/msg"
)

// process is the request dispatcher
func (h *NamespaceWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case msg.ActionAdd:
		h.add(q, &result)
	case msg.ActionRemove:
		h.remove(q, &result)
	case msg.ActionAttrAdd:
		h.attributeAdd(q, &result)
	case msg.ActionAttrRemove:
		h.attributeRemove(q, &result)
	case msg.ActionPropSet:
		h.propertySet(q, &result)
	case msg.ActionPropUpdate:
		h.propertyUpdate(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// remove deletes a specific namespace
func (h *NamespaceWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// propertyUpdate ...
func (h *NamespaceWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	var (
		err    error
		tx     *sql.Tx
		txTime time.Time
		ok     bool
	)
	txTime = time.Now().UTC()

	// tx.Begin()
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// forall Property in Request
	for key := range q.Namespace.Property {
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Namespace.Property[key],
		); !ok {
			tx.Rollback()
			return
		}
	}
	// tx.Commit()
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
