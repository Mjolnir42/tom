/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	NamespaceWriteStatements = ``

	NamespaceAdd = `
SELECT      'Namespace.ADD';`

	NamespaceRemove = `
SELECT      'Namespace.REMOVE';`
)

func init() {
	m[NamespaceAdd] = `NamespaceAdd`
	m[NamespaceRemove] = `NamespaceRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
