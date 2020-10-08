/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package server // import "github.com/mjolnir42/tom/internal/core/asset/server"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
)

// ReadHandler ...
type ReadHandler struct {
	Input         chan msg.Request
	Shutdown      chan struct{}
	name          string
	conn          *sql.DB
	lm            *lhm.LogHandleMap
	stmtAttribute *sql.Stmt
	stmtFind      *sql.Stmt
	stmtList      *sql.Stmt
	stmtParent    *sql.Stmt
	stmtLink      *sql.Stmt
}

func NewReadHandler(length int) (string, *ReadHandler) {
	h := &ReadHandler{}
	h.name = handler.GenerateName(`asset::server`) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
