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
	"github.com/mjolnir42/tom/pkg/proto"
)

// routeRegisterRuntime registers the runtime routes with the request
// router
func (m *Model) routeRegisterRuntime(rt *httprouter.Router) {

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

// handleRegisterRuntime registers the runtime application core handlers
// in the provided handlermap
func handleRegisterRuntime(hm *handler.Map, length int) {
	hm.Add(NewRuntimeReadHandler(length))
	hm.Add(NewRuntimeWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
