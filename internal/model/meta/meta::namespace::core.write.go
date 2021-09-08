/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// NamespaceWriteHandler is the handler for write requests changing
// namespace information
type NamespaceWriteHandler struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	name                 string
	conn                 *sql.DB
	lm                   *lhm.LogHandleMap
	stmtAdd              *sql.Stmt
	stmtConfig           *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtAttStdAdd        *sql.Stmt
	stmtAttUnqAdd        *sql.Stmt
	stmtAttDiscover      *sql.Stmt
	stmtAttQueryType     *sql.Stmt
	stmtTxStdPropAdd     *sql.Stmt
	stmtTxStdPropClamp   *sql.Stmt
	stmtTxStdPropSelect  *sql.Stmt
	stmtTxUniqPropAdd    *sql.Stmt
	stmtTxUniqPropClamp  *sql.Stmt
	stmtTxUniqPropSelect *sql.Stmt
}

// returns a new handler instance
func NewNamespaceWriteHandler(length int) (string, *NamespaceWriteHandler) {
	h := &NamespaceWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryMeta+`::`+msg.SectionNamespace) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *NamespaceWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionAttrAdd,
		proto.ActionAttrRemove,
		proto.ActionPropSet,
		proto.ActionPropUpdate,
		proto.ActionRemove,
	} {
		hm.Request(msg.SectionNamespace, action, h.name)
	}
}

// process is the request dispatcher
func (h *NamespaceWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionRemove:
		h.remove(q, &result)
	case proto.ActionAttrAdd:
		h.attributeAdd(q, &result)
	case proto.ActionAttrRemove:
		h.attributeRemove(q, &result)
	case proto.ActionPropSet:
		h.propertySet(q, &result)
	case proto.ActionPropUpdate:
		h.propertyUpdate(q, &result)
	case proto.ActionPropRemove:
		h.propertyRemove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *NamespaceWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *NamespaceWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *NamespaceWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for NamespaceWriteHandler
func (h *NamespaceWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceAdd:                  &h.stmtAdd,
		stmt.NamespaceConfigure:            &h.stmtConfig,
		stmt.NamespaceRemove:               &h.stmtRemove,
		stmt.NamespaceAttributeAddStandard: &h.stmtAttStdAdd,
		stmt.NamespaceAttributeAddUnique:   &h.stmtAttUnqAdd,
		stmt.NamespaceAttributeQueryType:   &h.stmtAttQueryType,
		stmt.NamespaceAttributeDiscover:    &h.stmtAttDiscover,
		stmt.NamespaceTxStdPropertyAdd:     &h.stmtTxStdPropAdd,
		stmt.NamespaceTxStdPropertyClamp:   &h.stmtTxStdPropClamp,
		stmt.NamespaceTxStdPropertySelect:  &h.stmtTxStdPropSelect,
		stmt.NamespaceTxUniqPropertyAdd:    &h.stmtTxUniqPropAdd,
		stmt.NamespaceTxUniqPropertyClamp:  &h.stmtTxUniqPropClamp,
		stmt.NamespaceTxUniqPropertySelect: &h.stmtTxUniqPropSelect,
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
func (h *NamespaceWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
