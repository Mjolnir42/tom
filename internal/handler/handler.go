/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package handler // import "github.com/mjolnir42/soma/internal/handler"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/msg"
)

// Handler process a specific request type
type Handler interface {
	Configure(*sql.DB, *lhm.LogHandleMap)
	Intake() chan msg.Request
	PriorityIntake() chan msg.Request
	Register(*Map)
	Run()
	ShutdownNow()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
