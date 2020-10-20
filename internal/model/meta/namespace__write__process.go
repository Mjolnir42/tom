/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	//"database/sql"

	//	"github.com/mjolnir42/lhm"
	//	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	//	"github.com/mjolnir42/tom/internal/stmt"
)

// process is the request dispatcher
func (h *NamespaceWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case msg.ActionAdd:
		h.add(q, &result)
	case msg.ActionRemove:
		h.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add creates a new namespace
func (h *NamespaceWriteHandler) add(q *msg.Request, mr *msg.Result) {
}

// remove deletes a specific namespace
func (h *NamespaceWriteHandler) remove(q *msg.Request, mr *msg.Result) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
