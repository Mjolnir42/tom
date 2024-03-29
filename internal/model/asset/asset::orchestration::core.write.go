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

// OrchestrationWriteHandler ...
type OrchestrationWriteHandler struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	name                 string
	conn                 *sql.DB
	lm                   *lhm.LogHandleMap
	stmtAdd              *sql.Stmt
	stmtAttAddStandard   *sql.Stmt
	stmtAttDiscover      *sql.Stmt
	stmtAttQueryType     *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtTxChild          *sql.Stmt
	stmtTxChildClean     *sql.Stmt
	stmtTxLink           *sql.Stmt
	stmtTxRteShow        *sql.Stmt
	stmtTxShow           *sql.Stmt
	stmtTxStackAdd       *sql.Stmt
	stmtTxStackClamp     *sql.Stmt
	stmtTxStackClampAll  *sql.Stmt
	stmtTxStdPropAdd     *sql.Stmt
	stmtTxStdPropClamp   *sql.Stmt
	stmtTxStdPropClean   *sql.Stmt
	stmtTxStdPropSelect  *sql.Stmt
	stmtTxUniqPropAdd    *sql.Stmt
	stmtTxUniqPropClamp  *sql.Stmt
	stmtTxUniqPropClean  *sql.Stmt
	stmtTxUniqPropSelect *sql.Stmt
}

// NewOrchestrationWriteHandler returns a new handler instance
func NewOrchestrationWriteHandler(length int) (string, *OrchestrationWriteHandler) {
	h := &OrchestrationWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryAsset+`::`+msg.SectionOrchestration) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *OrchestrationWriteHandler) Register(hm *handler.Map) {
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
		hm.Request(msg.SectionOrchestration, action, h.name)
	}
}

// process is the request dispatcher
func (h *OrchestrationWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

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
func (h *OrchestrationWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *OrchestrationWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *OrchestrationWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for OrchestrationWriteHandler
func (h *OrchestrationWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.NamespaceAttributeAddStandard:     &h.stmtAttAddStandard,
		stmt.NamespaceAttributeDiscover:        &h.stmtAttDiscover,
		stmt.NamespaceAttributeQueryType:       &h.stmtAttQueryType,
		stmt.OrchestrationAdd:                  &h.stmtAdd,
		stmt.OrchestrationTxLink:               &h.stmtTxLink,
		stmt.OrchestrationTxShow:               &h.stmtTxShow,
		stmt.OrchestrationTxStackAdd:           &h.stmtTxStackAdd,
		stmt.OrchestrationTxStackClamp:         &h.stmtTxStackClamp,
		stmt.OrchestrationTxStackClampAll:      &h.stmtTxStackClampAll,
		stmt.OrchestrationTxStdPropertyAdd:     &h.stmtTxStdPropAdd,
		stmt.OrchestrationTxStdPropertyClamp:   &h.stmtTxStdPropClamp,
		stmt.OrchestrationTxStdPropertyClean:   &h.stmtTxStdPropClean,
		stmt.OrchestrationTxStdPropertySelect:  &h.stmtTxStdPropSelect,
		stmt.OrchestrationTxUniqPropertyAdd:    &h.stmtTxUniqPropAdd,
		stmt.OrchestrationTxUniqPropertyClamp:  &h.stmtTxUniqPropClamp,
		stmt.OrchestrationTxUniqPropertyClean:  &h.stmtTxUniqPropClean,
		stmt.OrchestrationTxUniqPropertySelect: &h.stmtTxUniqPropSelect,
		stmt.OrchestrationTxUnstackChild:       &h.stmtTxChild,
		stmt.OrchestrationTxUnstackChildClean:  &h.stmtTxChildClean,
		stmt.RuntimeTxShow:                     &h.stmtTxRteShow,
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
func (h *OrchestrationWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
