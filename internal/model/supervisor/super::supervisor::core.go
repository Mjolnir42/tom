/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package supervisor // import "github.com/mjolnir42/tom/internal/model/supervisor/"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// SupervisorCoreHandler ...
type SupervisorCoreHandler struct {
	Input            chan msg.Request
	Shutdown         chan struct{}
	name             string
	conn             *sql.DB
	lm               *lhm.LogHandleMap
	stmtLinked       *sql.Stmt
	stmtList         *sql.Stmt
	stmtProp         *sql.Stmt
	stmtShow         *sql.Stmt
	stmtTxParent     *sql.Stmt
	stmtTxResolvNext *sql.Stmt
	stmtTxResolvPhys *sql.Stmt
	stmtTxResource   *sql.Stmt
}

// NewSupervisorCoreHandler returns a new handler instance
func NewSupervisorCoreHandler(length int) (string, *SupervisorCoreHandler) {
	h := &SupervisorCoreHandler{}
	h.name = handler.GenerateName(proto.ModelInternal+`::`+proto.EntitySupervisor) + `/core`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *SupervisorCoreHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAuthenticateEPK,
	} {
		hm.Request(proto.EntitySupervisor, action, h.name)
	}
}

// process is the request dispatcher
func (h *SupervisorCoreHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case proto.ActionAuthenticateEPK:
		h.authenticateEPK(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *SupervisorCoreHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *SupervisorCoreHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *SupervisorCoreHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for SupervisorCoreHandler
func (h *SupervisorCoreHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.ContainerList:              &h.stmtList,
		stmt.ContainerListLinked:        &h.stmtLinked,
		stmt.ContainerTxParent:          &h.stmtTxParent,
		stmt.ContainerTxResolvePhysical: &h.stmtTxResolvPhys,
		stmt.ContainerTxResolveServer:   &h.stmtTxResolvNext,
		stmt.ContainerTxSelectResource:  &h.stmtTxResource,
		stmt.ContainerTxShow:            &h.stmtShow,
		stmt.ContainerTxShowProperties:  &h.stmtProp,
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
func (h *SupervisorCoreHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
