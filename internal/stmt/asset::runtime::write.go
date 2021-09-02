/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	RuntimeWriteStatements = ``

	RuntimeAdd = `
SELECT      'Runtime.ADD';`

	RuntimeRemove = `
SELECT      'Runtime.REMOVE';`

	RuntimeStdAttrRemove = `
DELETE FROM       asset.runtime_environment_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	RuntimeUniqAttrRemove = `
DELETE FROM       asset.runtime_environment_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`
)

func init() {
	m[RuntimeAdd] = `RuntimeAdd`
	m[RuntimeRemove] = `RuntimeRemove`
	m[RuntimeStdAttrRemove] = `RuntimeStdAttrRemove`
	m[RuntimeUniqAttrRemove] = `RuntimeUniqAttrRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
