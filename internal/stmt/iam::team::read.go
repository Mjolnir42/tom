/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	TeamReadStatements = ``

	TeamList = `
SELECT      'TODO::TeamList'::text;`  // TODO

	TeamMbrList = `
SELECT      'TODO::TeamMbrList'::text;`  // TODO

	TeamShow = `
SELECT      'TODO::TeamShow'::text;`  // TODO
)

func init() {
	m[TeamList] = `TeamList`
	m[TeamMbrList] = `TeamMbrList`
	m[TeamShow] = `TeamShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
