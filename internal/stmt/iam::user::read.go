/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	UserReadStatements = ``

	UserList = `
SELECT      'TODO::UserList'::text;`  // TODO

	UserShow = `
SELECT      'TODO::UserShow'::text;`  // TODO
)

func init() {
	m[UserList] = `UserList`
	m[UserShow] = `UserShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
