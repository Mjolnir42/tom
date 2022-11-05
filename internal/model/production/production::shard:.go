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

// handleRegisterShard registers the blueprint application core handlers
// in the provided handlermap
func handleRegisterShard(hm *handler.Map, length int) {
	hm.Add(NewShardReadHandler(length))
	hm.Add(NewShardWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
