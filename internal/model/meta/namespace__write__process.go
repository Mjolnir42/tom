/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/model/meta"

import (
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
