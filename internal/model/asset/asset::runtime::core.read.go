/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// RuntimeReadHandler ...
type RuntimeReadHandler struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	name       string
	conn       *sql.DB
	lm         *lhm.LogHandleMap
	stmtLinked *sql.Stmt
	stmtList   *sql.Stmt
	stmtParent *sql.Stmt
	stmtProp   *sql.Stmt
	stmtShow   *sql.Stmt
	stmtTxChildren *sql.Stmt
}

// NewRuntimeReadHandler returns a new handler instance
func NewRuntimeReadHandler(length int) (string, *RuntimeReadHandler) {
	h := &RuntimeReadHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionRuntime) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *RuntimeReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionRuntime, action, h.name)
	}
}

// process is the request dispatcher
func (h *RuntimeReadHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case proto.ActionList:
		h.list(q, &result)
	case proto.ActionShow:
		h.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *RuntimeReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *RuntimeReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *RuntimeReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for RuntimeReadHandler
func (h *RuntimeReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.RuntimeList:             &h.stmtList,
		stmt.RuntimeListLinked:       &h.stmtLinked,
		stmt.RuntimeParent:           &h.stmtParent,
		stmt.RuntimeTxShow:           &h.stmtShow,
		stmt.RuntimeTxShowChildren:   &h.stmtTxChildren,
		stmt.RuntimeTxShowProperties: &h.stmtProp,
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
func (h *RuntimeReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
