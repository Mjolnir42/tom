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
	stmtLink             *sql.Stmt
	stmtTxCldCnr         *sql.Stmt
	stmtTxCldCnrClean    *sql.Stmt
	stmtTxCldOre         *sql.Stmt
	stmtTxCldOreClean    *sql.Stmt
	stmtTxCldRte         *sql.Stmt
	stmtTxCldRteClean    *sql.Stmt
	stmtTxCldSok         *sql.Stmt
	stmtTxCldSokClean    *sql.Stmt
	stmtTxCldSrv         *sql.Stmt
	stmtTxCldSrvClean    *sql.Stmt
	stmtTxOrchShow       *sql.Stmt
	stmtTxServerShow     *sql.Stmt
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
		proto.ActionLink,
		proto.ActionPropRemove,
		proto.ActionPropSet,
		proto.ActionPropUpdate,
		proto.ActionRemove,
		proto.ActionStack,
		proto.ActionUnstack,
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
		stmt.NamespaceAttributeQueryType:   &h.stmtAttQueryType,
		stmt.OrchestrationTxShow:           &h.stmtTxOrchShow,
		stmt.RuntimeAdd:                    &h.stmtAdd,
		stmt.RuntimeLink:                   &h.stmtLink,
		stmt.RuntimeTxShow:                 &h.stmtTxShow,
		stmt.RuntimeTxStackAdd:             &h.stmtTxStackAdd,
		stmt.RuntimeTxStackClamp:           &h.stmtTxStackClamp,
		stmt.RuntimeTxStdPropertyAdd:       &h.stmtTxStdPropAdd,
		stmt.RuntimeTxStdPropertyClamp:     &h.stmtTxStdPropClamp,
		stmt.RuntimeTxStdPropertyClean:     &h.stmtTxStdPropClean,
		stmt.RuntimeTxStdPropertySelect:    &h.stmtTxStdPropSelect,
		stmt.RuntimeTxUniqPropertyAdd:      &h.stmtTxUniqPropAdd,
		stmt.RuntimeTxUniqPropertyClamp:    &h.stmtTxUniqPropClamp,
		stmt.RuntimeTxUniqPropertyClean:    &h.stmtTxUniqPropClean,
		stmt.RuntimeTxUniqPropertySelect:   &h.stmtTxUniqPropSelect,
		stmt.RuntimeTxUnstackChildCnr:      &h.stmtTxCldCnr,
		stmt.RuntimeTxUnstackChildCnrClean: &h.stmtTxCldCnrClean,
		stmt.RuntimeTxUnstackChildOre:      &h.stmtTxCldOre,
		stmt.RuntimeTxUnstackChildOreClean: &h.stmtTxCldOreClean,
		stmt.RuntimeTxUnstackChildRte:      &h.stmtTxCldRte,
		stmt.RuntimeTxUnstackChildRteClean: &h.stmtTxCldRteClean,
		stmt.RuntimeTxUnstackChildSok:      &h.stmtTxCldSok,
		stmt.RuntimeTxUnstackChildSokClean: &h.stmtTxCldSokClean,
		stmt.RuntimeTxUnstackChildSrv:      &h.stmtTxCldSrv,
		stmt.RuntimeTxUnstackChildSrvClean: &h.stmtTxCldSrvClean,
		stmt.ServerTxShow:                  &h.stmtTxServerShow,
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
