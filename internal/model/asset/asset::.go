/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/rest"
)

var registry = make([]function, 0, 32)

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

func (m *Model) RouteRegister(rt *httprouter.Router) {
	m.routeRegisterServer(rt)
	m.routeRegisterRuntime(rt)
	m.routeRegisterOrchestration(rt)
}

// HandleRegister registers the application core handlers in the
// provided handlerMap
func HandleRegister(hm *handler.Map, length int) {
	handleRegisterRuntime(hm, length)
	handleRegisterServer(hm, length)
	handleRegisterOrchestration(hm, length)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
