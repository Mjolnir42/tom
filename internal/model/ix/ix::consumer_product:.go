//go:build exclude

/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ix // import "github.com/mjolnir42/tom/internal/model/ix/"

import (
	"github.com/mjolnir42/tom/internal/handler"
)

// handleRegisterConsumerProduct registers the blueprint application core handlers
// in the provided handlermap
func handleRegisterConsumerProduct(hm *handler.Map, length int) {
	hm.Add(NewConsumerProductReadHandler(length))
	hm.Add(NewConsumerProductWriteHandler(length))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
