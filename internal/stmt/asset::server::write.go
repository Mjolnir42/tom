/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	ServerWriteStatements = ``

	ServerAdd = `
SELECT      'Server.ADD';`

	ServerRemove = `
SELECT      'Server.REMOVE';`
)

func init() {
	m[ServerAdd] = `ServerAdd`
	m[ServerRemove] = `ServerRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
