/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package page // import "github.com/mjolnir42/tom/internal/model/page/"

import (
	"github.com/julienschmidt/httprouter"
)

//go:generate go-bindata -prefix bindata/ -pkg page bindata/... tmpl/...

// routeRegisterPage registers the page routes with the request router
func (m *Model) routeRegisterPage(rt *httprouter.Router) {
	// intentional no-op
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
