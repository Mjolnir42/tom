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

// ContainerWriteHandler ...
type ContainerWriteHandler struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	name                 string
	conn                 *sql.DB
	lm                   *lhm.LogHandleMap
	stmtAdd              *sql.Stmt
	stmtAttQueryType     *sql.Stmt
	stmtLink             *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtTxStdPropAdd     *sql.Stmt
	stmtTxStdPropClamp   *sql.Stmt
	stmtTxStdPropSelect  *sql.Stmt
	stmtTxUniqPropAdd    *sql.Stmt
	stmtTxUniqPropClamp  *sql.Stmt
	stmtTxUniqPropSelect *sql.Stmt
}

// NewContainerWriteHandler returns a new handler instance
func NewContainerWriteHandler(length int) (string, *ContainerWriteHandler) {
	h := &ContainerWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionContainer) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *ContainerWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionLink,
		proto.ActionPropRemove,
		proto.ActionPropSet,
		proto.ActionPropUpdate,
		proto.ActionRemove,
	} {
		hm.Request(msg.SectionContainer, action, h.name)
	}
}

// process is the request dispatcher
func (h *ContainerWriteHandler) process(q *msg.Request) {
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
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *ContainerWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *ContainerWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *ContainerWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for ContainerWriteHandler
func (h *ContainerWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceAttributeQueryType: &h.stmtAttQueryType,
		stmt.ContainerAdd:                  &h.stmtAdd,
		stmt.ContainerLink:                 &h.stmtLink,
		stmt.ContainerRemove:               &h.stmtRemove,
		stmt.ContainerTxStdPropertyAdd:     &h.stmtTxStdPropAdd,
		stmt.ContainerTxStdPropertyClamp:   &h.stmtTxStdPropClamp,
		stmt.ContainerTxStdPropertySelect:  &h.stmtTxStdPropSelect,
		stmt.ContainerTxUniqPropertyAdd:    &h.stmtTxUniqPropAdd,
		stmt.ContainerTxUniqPropertyClamp:  &h.stmtTxUniqPropClamp,
		stmt.ContainerTxUniqPropertySelect: &h.stmtTxUniqPropSelect,
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
func (h *ContainerWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
