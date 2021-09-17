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
// router for routes not added into the registry during init()
func (m *Model) routeRegisterRuntime(rt *httprouter.Router) {
}

// handleRegisterRuntime registers the runtime application core handlers
// in the provided handlermap
func handleRegisterRuntime(hm *handler.Map, length int) {
	hm.Add(NewRuntimeReadHandler(length))
	hm.Add(NewRuntimeWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
