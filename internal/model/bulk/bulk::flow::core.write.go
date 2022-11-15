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

// FlowWriteHandler ...
type FlowWriteHandler struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	name                 string
	conn                 *sql.DB
	lm                   *lhm.LogHandleMap
	stmtAdd              *sql.Stmt
	stmtAttQueryDiscover *sql.Stmt
	stmtAttQueryType     *sql.Stmt
	stmtTxShow           *sql.Stmt
	stmtTxStdPropAdd     *sql.Stmt
	stmtTxStdPropClamp   *sql.Stmt
	stmtTxStdPropClean   *sql.Stmt
	stmtTxStdPropSelect  *sql.Stmt
	stmtTxUniqPropAdd    *sql.Stmt
	stmtTxUniqPropClamp  *sql.Stmt
	stmtTxUniqPropClean  *sql.Stmt
	stmtTxUniqPropSelect *sql.Stmt
}

// NewFlowWriteHandler returns a new handler instance
func NewFlowWriteHandler(length int) (string, *FlowWriteHandler) {
	h := &FlowWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryBulk+`::`+msg.SectionFlow) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *FlowWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionEnsure,
		proto.ActionPropRemove,
		proto.ActionPropSet,
		proto.ActionPropUpdate,
		proto.ActionRemove,
	} {
		hm.Request(msg.SectionFlow, action, h.name)
	}
}

// process is the request dispatcher
func (h *FlowWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionEnsure:
		h.ensure(q, &result)
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
func (h *FlowWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *FlowWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *FlowWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for FlowWriteHandler
func (h *FlowWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.FlowAdd:                     &h.stmtAdd,
		stmt.FlowTxShow:                  &h.stmtTxShow,
		stmt.FlowTxStdPropertyAdd:        &h.stmtTxStdPropAdd,
		stmt.FlowTxStdPropertyClamp:      &h.stmtTxStdPropClamp,
		stmt.FlowTxStdPropertyClean:      &h.stmtTxStdPropClean,
		stmt.FlowTxStdPropertySelect:     &h.stmtTxStdPropSelect,
		stmt.FlowTxUniqPropertyAdd:       &h.stmtTxUniqPropAdd,
		stmt.FlowTxUniqPropertyClamp:     &h.stmtTxUniqPropClamp,
		stmt.FlowTxUniqPropertyClean:     &h.stmtTxUniqPropClean,
		stmt.FlowTxUniqPropertySelect:    &h.stmtTxUniqPropSelect,
		stmt.NamespaceAttributeDiscover:  &h.stmtAttQueryDiscover,
		stmt.NamespaceAttributeQueryType: &h.stmtAttQueryType,
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
func (h *FlowWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
