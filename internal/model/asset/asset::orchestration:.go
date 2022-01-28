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

// routeRegisterOrchestration registers the orchestration routes with the request
// router
func (m *Model) routeRegisterOrchestration(rt *httprouter.Router) {
}

// handleRegisterOrchestration registers the orchestration application core handlers
// in the provided handlermap
func handleRegisterOrchestration(hm *handler.Map, length int) {
	hm.Add(NewOrchestrationReadHandler(length))
	hm.Add(NewOrchestrationWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
