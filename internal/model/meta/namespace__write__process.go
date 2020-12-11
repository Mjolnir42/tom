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

// add creates a new namespace
func (h *NamespaceWriteHandler) add(q *msg.Request, mr *msg.Result) {
	// tx.Begin()
	// stmt.NamespaceAdd(proto.Namespace.Name)
	// stmt.NamespaceConfigure(proto.Namespace.Name, `dict_type`, proto.Namespace.Type
	// stmt.NamespaceConfigure(proto.Namespace.Name, `dict_lookup`, proto.Namespace.LookupKey)
	// stmt.NamespaceConfigure(proto.Namespace.Name, `dict_uri`, proto.Namespace.LookupURI)
	// stmt.NamespaceConfigure(proto.Namespace.Name, `dict_ntt_list`, // proto.Namespace.Constraint[])
	// tx.Commit()
}

// remove deletes a specific namespace
func (h *NamespaceWriteHandler) remove(q *msg.Request, mr *msg.Result) {
}

// attributeAdd ...
func (h *NamespaceWriteHandler) attributeAdd(q *msg.Request, mr *msg.Result) {
}

// attributeRemove ...
func (h *NamespaceWriteHandler) attributeRemove(q *msg.Request, mr *msg.Result) {
}

// propertySet ...
func (h *NamespaceWriteHandler) propertySet(q *msg.Request, mr *msg.Result) {
	// tx.Begin()
	// forall Property in Request
	// 		NamespaceTxStdPropertySelect -> already has value?
	// 		yes: NamespaceTxStdPropertyClamp
	// 		NamespaceTxStdPropertyAdd
	// forall Property currently set but not in Request
	//		NamespaceTxStdPropertyClamp
	// tx.Commit()
}

// propertyUpdate ...
func (h *NamespaceWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	// tx.Begin()
	// forall Property in Request
	// 		NamespaceTxStdPropertySelect -> already has value?
	// 		yes:NamespaceTxStdPropertyClamp
	// 		NamespaceTxStdPropertyAdd
	// tx.Commit()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix