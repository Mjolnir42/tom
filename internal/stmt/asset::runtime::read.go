/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	RuntimeReadStatements = ``

	RuntimeList = `
SELECT      'Runtime.LIST';`

	RuntimeShow = `
SELECT      'Runtime.SHOW';`
)

func init() {
	m[RuntimeList] = `RuntimeList`
	m[RuntimeShow] = `RuntimeShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
