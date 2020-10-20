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
SELECT      asset.server.serverID,
            meta.dictionary.name,
            meta.standard_attribute.attribute,
            asset.server_standard_attribute_values.value
FROM        asset.server
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.server_standard_attribute_values
      ON    asset.server_standard_attribute_values.serverID = asset.server.serverID
     AND    asset.server_standard_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_standard_attribute_values.attributeID = meta.standard_attribute.attributeID
   WHERE    now()::timestamptz(3) <@ asset.server_standard_attribute_values.validity
     AND    meta.standard_attribute.attribute IN ('type')
UNION
SELECT      asset.server.serverID,
            meta.dictionary.name,
            meta.unique_attribute.attribute,
            asset.server_unique_attribute_values.value
FROM        asset.server
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    asset.server_unique_attribute_values.serverID = asset.server.serverID
     AND    asset.server_unique_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
WHERE       now()::timestamptz(3) <@ asset.server_unique_attribute_values.validity
     AND    meta.unique_attribute.attribute IN ('name');`

	NamespaceShow = `
SELECT      'Namespace.SHOW';`
)

func init() {
	m[NamespaceList] = `NamespaceList`
	m[NamespaceShow] = `NamespaceShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
