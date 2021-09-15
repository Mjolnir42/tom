/*-
 * Copyright (c) 2020, Jörg Pernfuß
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

// RuntimeWriteHandler ...
type RuntimeWriteHandler struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	name                 string
	conn                 *sql.DB
	lm                   *lhm.LogHandleMap
	stmtAdd              *sql.Stmt
	stmtAttQueryType     *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtTxStdPropAdd     *sql.Stmt
	stmtTxStdPropClamp   *sql.Stmt
	stmtTxStdPropSelect  *sql.Stmt
	stmtTxUniqPropAdd    *sql.Stmt
	stmtTxUniqPropClamp  *sql.Stmt
	stmtTxUniqPropSelect *sql.Stmt
}

// NewRuntimeWriteHandler returns a new handler instance
func NewRuntimeWriteHandler(length int) (string, *RuntimeWriteHandler) {
	h := &RuntimeWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionRuntime) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *RuntimeWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionPropRemove,
		proto.ActionPropSet,
		proto.ActionPropUpdate,
		proto.ActionRemove,
	} {
		hm.Request(msg.SectionRuntime, action, h.name)
	}
}

// process is the request dispatcher
func (h *RuntimeWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionPropRemove:
		h.propertyRemove(q, &result)
	case proto.ActionPropSet:
		h.propertySet(q, &result)
	case proto.ActionPropUpdate:
		h.propertyUpdate(q, &result)
	case proto.ActionRemove:
		h.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *RuntimeWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *RuntimeWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *RuntimeWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for RuntimeWriteHandler
func (h *RuntimeWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceAttributeQueryType: &h.stmtAttQueryType,
		stmt.RuntimeTxStdPropertyAdd:     &h.stmtTxStdPropAdd,
		stmt.RuntimeTxStdPropertyClamp:   &h.stmtTxStdPropClamp,
		stmt.RuntimeTxStdPropertySelect:  &h.stmtTxStdPropSelect,
		stmt.RuntimeTxUniqPropertyAdd:    &h.stmtTxUniqPropAdd,
		stmt.RuntimeTxUniqPropertyClamp:  &h.stmtTxUniqPropClamp,
		stmt.RuntimeTxUniqPropertySelect: &h.stmtTxUniqPropSelect,
		stmt.RuntimeAdd:                  &h.stmtAdd,
		stmt.RuntimeRemove:               &h.stmtRemove,
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
func (h *RuntimeWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix