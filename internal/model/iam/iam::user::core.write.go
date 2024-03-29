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

// UserWriteHandler ...
type UserWriteHandler struct {
	Input         chan msg.Request
	Shutdown      chan struct{}
	name          string
	conn          *sql.DB
	lm            *lhm.LogHandleMap
	stmtActivate  *sql.Stmt
	stmtAdd       *sql.Stmt
	stmtAddKey    *sql.Stmt
	stmtDetect    *sql.Stmt
	stmtEnrolment *sql.Stmt
	stmtRemove    *sql.Stmt
	stmtUpdate    *sql.Stmt
	stmtUpdateUID *sql.Stmt
}

func NewUserWriteHandler(length int) (string, *UserWriteHandler) {
	h := &UserWriteHandler{}
	h.name = handler.GenerateName(msg.CategoryIAM+`::`+msg.SectionUser) + `/write`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// Register the handlername for the requests it wants to receive
func (h *UserWriteHandler) Register(hm *handler.Map) {
	for _, action := range []string{
		proto.ActionAdd,
		proto.ActionEnrolment,
		proto.ActionRemove,
		proto.ActionUpdate,
	} {
		hm.Request(msg.SectionUser, action, h.name)
	}
	for _, action := range []string{
		proto.ActionEnrolment,
	} {
		hm.Request(msg.SectionMachine, action, h.name)
	}
}

// process is the request dispatcher
func (h *UserWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Section {
	case msg.SectionUser:
		switch q.Action {
		case proto.ActionAdd:
			h.add(q, &result)
		case proto.ActionRemove:
			h.remove(q, &result)
		case proto.ActionUpdate:
			h.update(q, &result)
		case proto.ActionEnrolment:
			h.userEnrolment(q, &result)
		default:
			result.UnknownRequest(q)
		}
	case msg.SectionMachine:
		switch q.Action {
		case proto.ActionEnrolment:
			h.machineEnrolment(q, &result)
		default:
			result.UnknownRequest(q)
		}
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// Configure injects the handler with db connection and logging
func (h *UserWriteHandler) Configure(conn *sql.DB, lm *lhm.LogHandleMap) {
	h.conn = conn
	h.lm = lm
}

// Intake exposes the Input channel as part of the handler interface
func (h *UserWriteHandler) Intake() chan msg.Request {
	return h.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (h *UserWriteHandler) PriorityIntake() chan msg.Request {
	return h.Intake()
}

// Run is the event loop for UserWriteHandler
func (h *UserWriteHandler) Run() {
	var err error

	for statement, prepared := range map[string]**sql.Stmt{
		stmt.LibraryDetect:    &h.stmtDetect,
		stmt.MachineEnrol:     &h.stmtEnrolment,
		stmt.MachineUpdateUID: &h.stmtUpdateUID,
		stmt.UserActivate:     &h.stmtActivate,
		stmt.UserAdd:          &h.stmtAdd,
		stmt.UserAddKey:       &h.stmtAddKey,
		stmt.UserRemove:       &h.stmtRemove,
		stmt.UserUpdate:       &h.stmtUpdate,
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
func (h *UserWriteHandler) ShutdownNow() {
	close(h.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
