//go:build exclude

/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package production // import "github.com/mjolnir42/tom/internal/model/production/"

import (
	"github.com/mjolnir42/tom/internal/handler"
)

// handleRegisterDeployment registers the blueprint application core handlers
// in the provided handlermap
func handleRegisterDeployment(hm *handler.Map, length int) {
	hm.Add(NewDeploymentReadHandler(length))
	hm.Add(NewDeploymentWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
