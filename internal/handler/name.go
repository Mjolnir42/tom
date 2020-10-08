/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package handler // import "github.com/mjolnir42/tom/internal/handler/"

import (
	"github.com/satori/go.uuid"
)

// GenerateName returns randomized handler name with a provided prefix
func GenerateName(prefix string) string {
	return prefix + `/` + uuid.NewV4().String()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
