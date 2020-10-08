/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package stmt provides SQL statement string constants for TOM
package stmt // import "github.com/mjolnir42/tom/internal/stmt"

var m = make(map[string]string)

func Name(statement string) string {
	return m[statement]
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
