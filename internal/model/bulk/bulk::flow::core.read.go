/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/model/bulk/"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// FlowReadHandler ...
type FlowReadHandler struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	name     string
	conn     *sql.DB
	lm       *lhm.LogHandleMap
	stmtList *sql.Stmt
	stmtProp *sql.Stmt
	stmtShow *sql.Stmt
}

// NewFlowReadHandler returns a new handler instance
func NewFlowReadHandler(length int) (string, *FlowReadHandler) {
	h := &FlowReadHandler{}
	h.name = handler.GenerateName(msg.CategoryBulk+`::`+msg.SectionFlow) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *FlowReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionFlow, action, h.name)
	}
}

// process is the request dispatcher
func (h *FlowReadHandler) process(q *msg.Request) {
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
func (h *FlowReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *FlowReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *FlowReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for FlowReadHandler
func (h *FlowReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.FlowList:             &h.stmtList,
		stmt.FlowTxShow:           &h.stmtShow,
		stmt.FlowTxShowProperties: &h.stmtProp,
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
func (h *FlowReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
