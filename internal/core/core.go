/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package core implements the application handler core of Tom's service.
package core // import "github.com/mjolnir42/tom/internal/core/"

import (
	"database/sql"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/handler"
)

// Core is ther core application struct of Tom's service
type Core struct {
	hm   *handler.Map
	lm   *lhm.LogHandleMap
	db   *sql.DB
	conf *config.Configuration
}

// New returns a new application core
func New(
	hm *handler.Map,
	lm *lhm.LogHandleMap,
	db *sql.DB,
	conf *config.Configuration,
) *Core {
	x := Core{
		hm:   hm,
		lm:   lm,
		db:   db,
		conf: conf,
	}
	return &x
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
