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

// OrchestrationReadHandler ...
type OrchestrationReadHandler struct {
	Input            chan msg.Request
	Shutdown         chan struct{}
	name             string
	conn             *sql.DB
	lm               *lhm.LogHandleMap
	stmtList         *sql.Stmt
	stmtTxChildren   *sql.Stmt
	stmtTxLinks      *sql.Stmt
	stmtTxParent     *sql.Stmt
	stmtTxProp       *sql.Stmt
	stmtTxResolvNext *sql.Stmt
	stmtTxResolvPhys *sql.Stmt
	stmtTxResource   *sql.Stmt
	stmtTxShow       *sql.Stmt
}

// NewOrchestrationReadHandler returns a new handler instance
func NewOrchestrationReadHandler(length int) (string, *OrchestrationReadHandler) {
	h := &OrchestrationReadHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionOrchestration) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *OrchestrationReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionResolve,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionOrchestration, action, h.name)
	}
}

// process is the request dispatcher
func (h *OrchestrationReadHandler) process(q *msg.Request) {
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
func (h *OrchestrationReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *OrchestrationReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *OrchestrationReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for OrchestrationReadHandler
func (h *OrchestrationReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.OrchestrationList:              &h.stmtList,
		stmt.OrchestrationListLinked:        &h.stmtTxLinks,
		stmt.OrchestrationTxParent:          &h.stmtTxParent,
		stmt.OrchestrationTxResolvePhysical: &h.stmtTxResolvPhys,
		stmt.OrchestrationTxResolveServer:   &h.stmtTxResolvNext,
		stmt.OrchestrationTxSelectResource:  &h.stmtTxResource,
		stmt.OrchestrationTxShow:            &h.stmtTxShow,
		stmt.OrchestrationTxShowChildren:    &h.stmtTxChildren,
		stmt.OrchestrationTxShowProperties:  &h.stmtTxProp,
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
func (h *OrchestrationReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
