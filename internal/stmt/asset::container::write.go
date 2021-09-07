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

	ContainerNamespaceRemoveLinking = `
DELETE FROM       asset.container_linking
WHERE             dictionaryID_A = $1::uuid
   OR             dictionaryID_B = $1::uuid;`

	ContainerNamespaceRemoveParent = `
DELETE FROM       asset.container_parent
USING             asset.container
WHERE             asset.container_parent.containerID = asset.container.containerID
  AND             asset.container.dictionaryID = $1::uuid;`

	ContainerNamespaceRemove = `
DELETE FROM       asset.container
WHERE             dictionaryID = $1::uuid;`
)

func init() {
	m[ContainerNamespaceRemoveLinking] = `ContainerNamespaceRemoveLinking`
	m[ContainerNamespaceRemoveParent] = `ContainerNamespaceRemoveParent`
	m[ContainerNamespaceRemove] = `ContainerNamespaceRemove`
	m[ContainerStdAttrRemove] = `ContainerStdAttrRemove`
	m[ContainerUniqAttrRemove] = `ContainerUniqAttrRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
