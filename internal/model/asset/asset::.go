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
	"github.com/mjolnir42/tom/pkg/proto"
)

var registry = make([]function, 0, 32)

type function struct {
	cmd    string
	handle func(*Model) httprouter.Handle
}

type Model struct {
	x *rest.Rest
}

// New returns a new instance of asset.Model
func New(x *rest.Rest) *Model {
	m := &Model{
		x: x,
	}
	return m
}

// RouteRegister registers all request routes for the asset package with
// the provided httprouter.Router instance
func (m *Model) RouteRegister(rt *httprouter.Router) {
	// these functions are currently no-op but left in place as
	// on-demand registration hook
	m.routeRegisterContainer(rt)
	m.routeRegisterOrchestration(rt)
	m.routeRegisterRuntime(rt)
	m.routeRegisterServer(rt)
	//m.routeRegisterSocket(rt)

	m.routeRegisterFromRegistry(rt)
}

// routeRegisterFromRegistry picks up the route registration requests
// put into the registry package variable during init() phase of the
// package
func (m *Model) routeRegisterFromRegistry(rt *httprouter.Router) {
	for _, f := range registry {
		m.x.LM.GetLogger(`application`).Infof(
			"Registering handle for %s at route %s|%s",
			f.cmd,
			proto.Commands[f.cmd].Method,
			proto.Commands[f.cmd].Path,
		)
		switch proto.Commands[f.cmd].Method {
		case proto.MethodDELETE:
			rt.DELETE(proto.Commands[f.cmd].Path, f.handle(m))
		case proto.MethodGET:
			rt.GET(proto.Commands[f.cmd].Path, f.handle(m))
		case proto.MethodHEAD:
			rt.HEAD(proto.Commands[f.cmd].Path, f.handle(m))
		case proto.MethodPATCH:
			rt.PATCH(proto.Commands[f.cmd].Path, f.handle(m))
		case proto.MethodPOST:
			rt.POST(proto.Commands[f.cmd].Path, f.handle(m))
		case proto.MethodPUT:
			rt.PUT(proto.Commands[f.cmd].Path, f.handle(m))
		default:
			m.x.LM.GetLogger(`error`).Errorf(
				"Error registering route for %s using unknown method %s",
				f.cmd,
				proto.Commands[f.cmd].Method,
			)
		}
	}
}

// HandleRegister registers the application core handlers in the
// provided handlerMap
func HandleRegister(hm *handler.Map, length int) {
	handleRegisterContainer(hm, length)
	handleRegisterOrchestration(hm, length)
	handleRegisterRuntime(hm, length)
	handleRegisterServer(hm, length)
	//handleRegisterSocket(hm, length)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
