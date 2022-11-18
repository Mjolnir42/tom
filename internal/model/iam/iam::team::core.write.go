/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// TeamWriteHandler ...
type TeamWriteHandler struct {
	Input         chan msg.Request
	Shutdown      chan struct{}
	name          string
	conn          *sql.DB
	lm            *lhm.LogHandleMap
	stmtAdd       *sql.Stmt
	stmtRemove    *sql.Stmt
	stmtUpdate    *sql.Stmt
	stmtHdSet     *sql.Stmt
	stmtHdUnset   *sql.Stmt
	stmtMbrAdd    *sql.Stmt
	stmtMbrSet    *sql.Stmt
	stmtMbrRemove *sql.Stmt
}

func NewTeamWriteHandler(length int) (string, *TeamWriteHandler) {
	h := &TeamWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryIAM+`::`+msg.SectionTeam) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *TeamWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionRemove,
		proto.ActionUpdate,
		proto.ActionHdSet,
		proto.ActionHdUnset,
		proto.ActionMbrAdd,
		proto.ActionMbrSet,
		proto.ActionMbrRemove,
	} {
		hm.Request(msg.SectionTeam, action, h.name)
	}
}

// process is the request dispatcher
func (h *TeamWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case proto.ActionAdd:
		h.add(q, &result)
	case proto.ActionRemove:
		h.remove(q, &result)
	case proto.ActionUpdate:
		h.update(q, &result)
	case proto.ActionHdSet:
		h.headOfSet(q, &result)
	case proto.ActionHdUnset:
		h.headOfUnset(q, &result)
	case proto.ActionMbrAdd:
		h.memberAdd(q, &result)
	case proto.ActionMbrSet:
		h.memberSet(q, &result)
	case proto.ActionMbrRemove:
		h.memberRemove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *TeamWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *TeamWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *TeamWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for TeamWriteHandler
func (h *TeamWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.TeamAdd:       &h.stmtAdd,
		stmt.TeamRemove:    &h.stmtRemove,
		stmt.TeamUpdate:    &h.stmtUpdate,
		stmt.TeamHdSet:     &h.stmtHdSet,
		stmt.TeamHdUnset:   &h.stmtHdUnset,
		stmt.TeamMbrAdd:    &h.stmtMbrAdd,
		stmt.TeamMbrSet:    &h.stmtMbrSet,
		stmt.TeamMbrRemove: &h.stmtMbrRemove,
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
func (h *TeamWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
