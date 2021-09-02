/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	ContainerWriteStatements = ``

	ContainerStdAttrRemove = `
DELETE FROM       asset.container_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	ContainerUniqAttrRemove = `
DELETE FROM       asset.container_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`
)

func init() {
	m[ContainerStdAttrRemove] = `ContainerStdAttrRemove`
	m[ContainerUniqAttrRemove] = `ContainerUniqAttrRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
