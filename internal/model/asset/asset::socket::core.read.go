// +build socket

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

// SocketReadHandler ...
type SocketReadHandler struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	name       string
	conn       *sql.DB
	lm         *lhm.LogHandleMap
	stmtLinked *sql.Stmt
	stmtList   *sql.Stmt
	stmtProp   *sql.Stmt
	stmtShow   *sql.Stmt
}

// NewSocketReadHandler returns a new handler instance
func NewSocketReadHandler(length int) (string, *SocketReadHandler) {
	h := &SocketReadHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionSocket) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *SocketReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionSocket, action, h.name)
	}
}

// process is the request dispatcher
func (h *SocketReadHandler) process(q *msg.Request) {
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
func (h *SocketReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *SocketReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *SocketReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for SocketReadHandler
func (h *SocketReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.SocketList:             &h.stmtList,
		stmt.SocketListLinked:       &h.stmtLinked,
		stmt.SocketTxShow:           &h.stmtShow,
		stmt.SocketTxShowProperties: &h.stmtProp,
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
func (h *SocketReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
