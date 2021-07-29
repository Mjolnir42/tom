/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	UserWriteStatements = ``

	UserAdd = `
SELECT      'TODO::UserAdd'::text;`  // TODO

	UserRemove = `
SELECT      'TODO::UserRemove'::text;`  // TODO

	UserUpdate = `
SELECT      'TODO::UserUpdate'::text;`  // TODO
)

func init() {
	m[UserAdd] = `UserAdd`
	m[UserRemove] = `UserRemove`
	m[UserUpdate] = `UserUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
