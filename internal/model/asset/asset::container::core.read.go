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

// ContainerReadHandler ...
type ContainerReadHandler struct {
	Input          chan msg.Request
	Shutdown       chan struct{}
	name           string
	conn           *sql.DB
	lm             *lhm.LogHandleMap
	stmtLinked     *sql.Stmt
	stmtList       *sql.Stmt
	stmtProp       *sql.Stmt
	stmtResolvNext *sql.Stmt
	stmtResolvPhys *sql.Stmt
	stmtShow       *sql.Stmt
	stmtTxParent   *sql.Stmt
	stmtTxResource *sql.Stmt
}

// NewContainerReadHandler returns a new handler instance
func NewContainerReadHandler(length int) (string, *ContainerReadHandler) {
	h := &ContainerReadHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionContainer) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *ContainerReadHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionList,
		proto.ActionResolve,
		proto.ActionShow,
	} {
		hm.Request(msg.SectionContainer, action, h.name)
	}
}

// process is the request dispatcher
func (h *ContainerReadHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case proto.ActionList:
		h.list(q, &result)
	case proto.ActionShow:
		h.show(q, &result)
	case proto.ActionResolve:
		h.resolve(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *ContainerReadHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *ContainerReadHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *ContainerReadHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for ContainerReadHandler
func (h *ContainerReadHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.ContainerList:             &h.stmtList,
		stmt.ContainerListLinked:       &h.stmtLinked,
		stmt.ContainerResolvePhysical:  &h.stmtResolvPhys,
		stmt.ContainerResolveServer:    &h.stmtResolvNext,
		stmt.ContainerTxParent:         &h.stmtTxParent,
		stmt.ContainerTxSelectResource: &h.stmtTxResource,
		stmt.ContainerTxShow:           &h.stmtShow,
		stmt.ContainerTxShowProperties: &h.stmtProp,
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
func (h *ContainerReadHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
