/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
)

// Implementation of the handler.Handler interface

// Configure injects the handler with db connection and logging
func (h *NamespaceWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *NamespaceWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *NamespaceWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Register the handlername for the requests it wants to receive
func (h *NamespaceWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionShow,
		msg.ActionAttrAdd,
		msg.ActionAttrRemove,
		msg.ActionPropSet,
		msg.ActionPropUpdate,
	} {
		hm.Request(msg.SectionNamespace, action, h.name)
	}
}

// Run is the event loop for NamespaceWriteHandler
func (h *NamespaceWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceAdd:    &h.stmtAdd,
		stmt.NamespaceRemove: &h.stmtRemove,
		stmt.NamespaceConfigure: &h.stmtConfig,
		stmt.NamespaceAttributeAddStandard: &h.stmtAttStdAdd,
		stmt.NamespaceAttributeAddUnique:   &h.stmtAddUnqAdd,
		stmt.NamespaceTxStdPropertyAdd:     &h.stmtTxStdPropAdd,
		stmt.NamespaceTxStdPropertyClamp:   &h.stmtTxStdPropClamp,
		stmt.NamespaceTxStdPropertySelect:  &h.stmtTxStdPropSelect,
		stmt.NamespaceTxUniqPropertyAdd:    &h.stmtTxUniqPropAdd,
		stmt.NamespaceTxUniqPropertyClamp:  &h.stmtTxUniqPropClamp,
		stmt.NamespaceTxUniqPropertySelect: &h.stmtTxUniqPropSelect,
	} {
		if *prepared, err = h.conn.Prepare(statement); err != nil {
			h.lm.GetLogger(`error`).Fatal(h.name, err, stmt.Name(statement))
		}
		defer (*prepared).Close()
	}

	for {
		select {
		case <-h.Shutdown:
			break
		case req := <-h.Input:
			go func() {
				h.process(&req)
			}()
		}
	}
}

// ShutdownNow signals the handler to shut down
func (h *NamespaceWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
