/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/model/bulk/"

import (
	"github.com/mjolnir42/tom/internal/handler"
)

// handleRegisterFlow registers the flow application core handlers
// in the provided handlermap
func handleRegisterFlow(hm *handler.Map, length int) {
	hm.Add(NewFlowReadHandler(length))
	hm.Add(NewFlowWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
