/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

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
func (h *UserWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *UserWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *UserWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Register the handlername for the requests it wants to receive
func (h *UserWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionRemove,
		proto.ActionUpdate,
	} {
		hm.Request(msg.SectionUser, action, h.name)
	}
}

// Run is the event loop for UserWriteHandler
func (h *UserWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.UserAdd:    &h.stmtAdd,
		stmt.UserRemove: &h.stmtRemove,
		stmt.UserUpdate: &h.stmtUpdate,
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
func (h *UserWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
