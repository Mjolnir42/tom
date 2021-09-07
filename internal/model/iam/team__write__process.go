/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

// process is the request dispatcher
func (h *TeamWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionRemove:
		h.remove(q, &result)
	case proto.ActionUpdate:
		h.update(q, &result)
	case proto.ActionHdSet:
		h.headOfSet(q, &result)
	case proto.ActionHdUnset:
		h.headOfUnset(q, &result)
	case proto.ActionMbrAdd:
		h.memberAdd(q, &result)
	case proto.ActionMbrRemove:
		h.memberRemove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add ...
func (h *TeamWriteHandler) add(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// remove ...
func (h *TeamWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// update ...
func (h *TeamWriteHandler) update(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// headOfSet ...
func (h *TeamWriteHandler) headOfSet(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// headOfUnset ...
func (h *TeamWriteHandler) headOfUnset(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// memberAdd ...
func (h *TeamWriteHandler) memberAdd(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// memberRemove ...
func (h *TeamWriteHandler) memberRemove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
