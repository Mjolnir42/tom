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

// ServerWriteHandler ...
type ServerWriteHandler struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	name                 string
	conn                 *sql.DB
	lm                   *lhm.LogHandleMap
	stmtAdd              *sql.Stmt
	stmtAttQueryType     *sql.Stmt
	stmtLink             *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtTxRuntimeShow    *sql.Stmt
	stmtTxShow           *sql.Stmt
	stmtTxStackAdd       *sql.Stmt
	stmtTxStackClamp     *sql.Stmt
	stmtTxStdPropAdd     *sql.Stmt
	stmtTxStdPropClamp   *sql.Stmt
	stmtTxStdPropClean   *sql.Stmt
	stmtTxStdPropSelect  *sql.Stmt
	stmtTxUniqPropAdd    *sql.Stmt
	stmtTxUniqPropClamp  *sql.Stmt
	stmtTxUniqPropClean  *sql.Stmt
	stmtTxUniqPropSelect *sql.Stmt
}

// NewServerWriteHandler returns a new handler instance
func NewServerWriteHandler(length int) (string, *ServerWriteHandler) {
	h := &ServerWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionServer) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *ServerWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionLink,
		proto.ActionPropRemove,
		proto.ActionPropSet,
		proto.ActionPropUpdate,
		proto.ActionRemove,
		proto.ActionStack,
		proto.ActionUnstack,
	} {
		hm.Request(msg.SectionServer, action, h.name)
	}
}

// process is the request dispatcher
func (h *ServerWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionLink:
		h.link(q, &result)
	case proto.ActionPropRemove:
		h.propertyRemove(q, &result)
	case proto.ActionPropSet:
		h.propertySet(q, &result)
	case proto.ActionPropUpdate:
		h.propertyUpdate(q, &result)
	case proto.ActionRemove:
		h.remove(q, &result)
	case proto.ActionStack:
		h.stack(q, &result)
	case proto.ActionUnstack:
		h.unstack(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *ServerWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *ServerWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *ServerWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for ServerWriteHandler
func (h *ServerWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceAttributeQueryType: &h.stmtAttQueryType,
		stmt.RuntimeTxShow:               &h.stmtTxRuntimeShow,
		stmt.ServerAdd:                   &h.stmtAdd,
		stmt.ServerLink:                  &h.stmtLink,
		stmt.ServerRemove:                &h.stmtRemove,
		stmt.ServerTxShow:                &h.stmtTxShow,
		stmt.ServerTxStackAdd:            &h.stmtTxStackAdd,
		stmt.ServerTxStackClamp:          &h.stmtTxStackClamp,
		stmt.ServerTxStdPropertyAdd:      &h.stmtTxStdPropAdd,
		stmt.ServerTxStdPropertyClamp:    &h.stmtTxStdPropClamp,
		stmt.ServerTxStdPropertyClean:    &h.stmtTxStdPropClean,
		stmt.ServerTxStdPropertySelect:   &h.stmtTxStdPropSelect,
		stmt.ServerTxUniqPropertyAdd:     &h.stmtTxUniqPropAdd,
		stmt.ServerTxUniqPropertyClamp:   &h.stmtTxUniqPropClamp,
		stmt.ServerTxUniqPropertyClean:   &h.stmtTxUniqPropClean,
		stmt.ServerTxUniqPropertySelect:  &h.stmtTxUniqPropSelect,
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
func (h *ServerWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
