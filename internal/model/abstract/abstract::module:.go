/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package abstract // import "github.com/mjolnir42/tom/internal/model/abstract/"

import (
	"github.com/mjolnir42/tom/internal/handler"
)

// handleRegisterModule registers the blueprint application core handlers
// in the provided handlermap
func handleRegisterModule(hm *handler.Map, length int) {
	hm.Add(NewModuleReadHandler(length))
	hm.Add(NewModuleWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
