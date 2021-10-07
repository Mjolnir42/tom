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
)

// routeRegisterSocket registers the socket routes with the request
// router for routes not added into the registry during init()
func (m *Model) routeRegisterSocket(rt *httprouter.Router) {
}

// handleRegisterSocket registers the socket application core handlers
// in the provided handlermap
func handleRegisterSocket(hm *handler.Map, length int) {
	//hm.Add(NewSocketReadHandler(length))
	//hm.Add(NewSocketWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
