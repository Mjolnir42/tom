/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/handler"
)

// routeRegisterLibrary registers the library routes with the
// request router
func (m *Model) routeRegisterLibrary(rt *httprouter.Router) {
	rt.GET(`/idlib/`, m.x.Authenticated(m.LibraryList))
	rt.GET(`/idlib/:lib`, m.x.Authenticated(m.LibraryShow))
	rt.POST(`/idlib/`, m.x.Authenticated(m.LibraryAdd))
	rt.DELETE(`/idlib/:lib`, m.x.Authenticated(m.LibraryRemove))
}

// handleRegisterLibrary registers library application core handlers in
// the provided handlermap
func handleRegisterLibrary(hm *handler.Map, length int) {
	hm.Add(NewLibraryReadHandler(length))
	hm.Add(NewLibraryWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
