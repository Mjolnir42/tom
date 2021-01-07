/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	NamespaceReadStatements = ``

	NamespaceList = `
SELECT      meta.dictionary.name,
            meta.dictionary.createdAt,
            inventory.user.uid
FROM        meta.dictionary
    JOIN    inventory.user
      ON    meta.dictionary.createdBy = inventory.user.userID;`

	NamespaceShow = `
SELECT      'Namespace.SHOW';`
)

func init() {
	m[NamespaceList] = `NamespaceList`
	m[NamespaceShow] = `NamespaceShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
