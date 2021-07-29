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
)

// UserReadHandler ...
type UserReadHandler struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	name     string
	conn     *sql.DB
	lm       *lhm.LogHandleMap
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
}

func NewUserReadHandler(length int) (string, *UserReadHandler) {
	h := &UserReadHandler{}
	h.name = handler.GenerateName(msg.CategoryIAM+`::`+msg.SectionUser) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
