/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/rest"
)

var registry = make([]function, 0, 12)

type function struct {
	cmd    string
	handle func(*Model) httprouter.Handle
}

type Model struct {
	x *rest.Rest
}

func New(x *rest.Rest) *Model {
	m := &Model{
		x: x,
	}
	return m
}

// RouteRegister registers all REST routes with the provided request
// router
func (m *Model) RouteRegister(rt *httprouter.Router) {
	m.routeRegisterNamespace(rt)
}

// HandleRegister registers the application core handlers in the
// provided handlerMap
func HandleRegister(hm *handler.Map, length int) {
	handleRegisterNamespace(hm, length)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
