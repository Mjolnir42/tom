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

// LibraryWriteHandler ...
type LibraryWriteHandler struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	name       string
	conn       *sql.DB
	lm         *lhm.LogHandleMap
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
}

func NewLibraryWriteHandler(length int) (string, *LibraryWriteHandler) {
	h := &LibraryWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryIAM+`::`+msg.SectionLibrary) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *LibraryWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionRemove,
	} {
		hm.Request(msg.SectionLibrary, action, h.name)
	}
}

// process is the request dispatcher
func (h *LibraryWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

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

// Configure injects the handler with db connection and logging
func (h *LibraryWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *LibraryWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *LibraryWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for LibraryWriteHandler
func (h *LibraryWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.LibraryAdd:    &h.stmtAdd,
		stmt.LibraryRemove: &h.stmtRemove,
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
func (h *LibraryWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
