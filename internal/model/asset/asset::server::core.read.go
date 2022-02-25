/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
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

// ServerReadHandler ...
type ServerReadHandler struct {
	Input            chan msg.Request
	Shutdown         chan struct{}
	name             string
	conn             *sql.DB
	lm               *lhm.LogHandleMap
	stmtAttribute    *sql.Stmt
	stmtFind         *sql.Stmt
	stmtLinked       *sql.Stmt
	stmtList         *sql.Stmt
	stmtParent       *sql.Stmt
	stmtTxChildren   *sql.Stmt
	stmtTxProp       *sql.Stmt
	stmtTxResolvNext *sql.Stmt
	stmtTxResolvPhys *sql.Stmt
	stmtTxResource   *sql.Stmt
	stmtTxShow       *sql.Stmt
}

// NewServerReadHandler returns a new handler instance
func NewServerReadHandler(length int) (string, *ServerReadHandler) {
	h := &ServerReadHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionServer) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *ServerReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionResolve,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionServer, action, h.name)
	}
}

// process is the request dispatcher
func (h *ServerReadHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case proto.ActionList:
		h.list(q, &result)
	case proto.ActionResolve:
		h.resolve(q, &result)
	case proto.ActionShow:
		h.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *ServerReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *ServerReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *ServerReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for ServerReadHandler
func (h *ServerReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.ServerAttribute:         &h.stmtAttribute,
		stmt.ServerFind:              &h.stmtFind,
		stmt.ServerList:              &h.stmtList,
		stmt.ServerListLinked:        &h.stmtLinked,
		stmt.ServerParent:            &h.stmtParent,
		stmt.ServerTxResolvePhysical: &h.stmtTxResolvPhys,
		stmt.ServerTxResolveServer:   &h.stmtTxResolvNext,
		stmt.ServerTxSelectResource:  &h.stmtTxResource,
		stmt.ServerTxShow:            &h.stmtTxShow,
		stmt.ServerTxShowChildren:    &h.stmtTxChildren,
		stmt.ServerTxShowProperties:  &h.stmtTxProp,
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
func (h *ServerReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
