/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"github.com/mjolnir42/tom/internal/msg"
)

// process is the request dispatcher
func (h *UserReadHandler) process(q *msg.Request) {
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
func (h *UserReadHandler) list(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// show returns full details for a specific server
func (h *UserReadHandler) show(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
