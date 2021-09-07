/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
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
	"github.com/mjolnir42/tom/pkg/proto"
)

// Implementation of the handler.Handler interface

// Configure injects the handler with db connection and logging
func (h *NamespaceReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *NamespaceReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *NamespaceReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Register the handlername for the requests it wants to receive
func (h *NamespaceReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionNamespace, action, h.name)
	}
}

// Run is the event loop for NamespaceReadHandler
func (h *NamespaceReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceList:               &h.stmtList,
		stmt.NamespaceTxSelectAttributes: &h.stmtAttr,
		stmt.NamespaceTxSelectProperties: &h.stmtProp,
		stmt.NamespaceTxShow:             &h.stmtShow,
	} {
		if *prepared, err = h.conn.Prepare(statement); err != nil {
			h.lm.GetLogger(`error`).Fatal(handler.StmtErr(h.name, err, stmt.Name(statement)))
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
func (h *NamespaceReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
