/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/tom/internal/model/super/"

import (
	"github.com/mjolnir42/tom/internal/handler"
)

// handleRegisterSupervisor registers the supervisor application core handlers
// in the provided handlermap
func handleRegisterSupervisor(hm *handler.Map, length int) {
	hm.Add(NewSupervisorCoreHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
