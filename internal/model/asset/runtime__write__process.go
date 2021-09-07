/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	//	"database/sql"
	//	"fmt"

	"github.com/mjolnir42/tom/internal/msg"
	//	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// process is the request dispatcher
func (h *RuntimeWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionRemove:
		h.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add ...
func (h *RuntimeWriteHandler) add(q *msg.Request, mr *msg.Result) {
}

// remove ...
func (h *RuntimeWriteHandler) remove(q *msg.Request, mr *msg.Result) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
