/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
)

// NamespaceReadHandler ...
type NamespaceReadHandler struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	name     string
	conn     *sql.DB
	lm       *lhm.LogHandleMap
	stmtAttr *sql.Stmt
	stmtList *sql.Stmt
	stmtProp *sql.Stmt
	stmtShow *sql.Stmt
}

func NewNamespaceReadHandler(length int) (string, *NamespaceReadHandler) {
	h := &NamespaceReadHandler{}
	h.name = handler.GenerateName(msg.CategoryMeta+`::`+msg.SectionNamespace) + `/read`
	h.Input = make(chan msg.Request, length)
	h.Shutdown = make(chan struct{})
	return h.name, h
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
