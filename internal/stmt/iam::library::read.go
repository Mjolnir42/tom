/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	LibraryReadStatements = ``

	LibraryList = `
SELECT      'TODO::LibraryList'::text;`  // TODO

	LibraryShow = `
SELECT      'TODO::LibraryShow'::text;`  // TODO
)

func init() {
	m[LibraryList] = `LibraryList`
	m[LibraryShow] = `LibraryShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
