/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"net/http"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
)

type Rest struct {
	isAuthorized func(*msg.Request) bool
	conf         *config.Configuration
	hm           *handler.Map
	lm           *lhm.LogHandleMap
	idx          int
}

func New(
	authorizationFunction func(*msg.Request) bool,
	index int,
	hm *handler.Map,
	lm *lhm.LogHandleMap,
	cfg *config.Configuration,
) *Rest {
	x := Rest{}
	x.isAuthorized = authorizationFunction
	x.idx = index
	x.hm = hm
	x.lm = lm
	x.conf = cfg
	return &x
}

func (x *Rest) Run() {
	router := x.setupRouter()

	// TODO switch to new abortable interface
	x.lm.GetLogger(`error`).Fatal(http.ListenAndServeTLS(
		x.conf.Daemon[x.idx].URL.Host,
		x.conf.Daemon[x.idx].Cert,
		x.conf.Daemon[x.idx].Key,
		router,
	))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
