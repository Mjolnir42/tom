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
	rt.GET(`/orchestration/`, m.x.Authenticated(m.OrchestrationList))
	rt.GET(`/orchestration/:tomID`, m.x.Authenticated(m.OrchestrationShow))
	rt.POST(`/orchestration/`, m.x.Authenticated(m.OrchestrationAdd))
	rt.DELETE(`/orchestration/:tomID`, m.x.Authenticated(m.OrchestrationRemove))
}

// handleRegisterOrchestration registers the orchestration application core handlers
// in the provided handlermap
func handleRegisterOrchestration(hm *handler.Map, length int) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
