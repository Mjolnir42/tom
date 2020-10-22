/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	ServerReadStatements = ``

	ServerAttribute = `
SELECT      asset.server.serverID,
            meta.dictionary.name,
            meta.standard_attribute.attribute,
            asset.server_standard_attribute_values.value,
			'standard'::text AS type
FROM        asset.server
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.server_standard_attribute_values
      ON    asset.server_standard_attribute_values.serverID = asset.server.serverID
     AND    asset.server_standard_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_standard_attribute_values.attributeID = meta.standard_attribute.attributeID
WHERE       now()::timestamptz(3) <@ asset.server_standard_attribute_values.validity
     AND    asset.server.serverID = $1::uuid
UNION
SELECT      asset.server.serverID,
            meta.dictionary.name,
            meta.unique_attribute.attribute,
            asset.server_unique_attribute_values.value,
			'unique'::text AS type
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
     AND    asset.server.serverID = $1::uuid;`

	ServerFind = `
SELECT      asset.server.serverID,
            asset.server.dictionaryID,
            meta.dictionary.name,
            meta.unique_attribute.attributeID,
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
     AND    meta.unique_attribute.attribute IN ('name')
     AND    (asset.server_unique_attribute_values.value = $1::varchar OR $1::varchar IS NULL)
     AND    (asset.server.serverID = $2::uuid OR $2::uuid IS NULL)
     AND    (asset.server.dictionaryID = $3::uuid OR $3::uuid IS NULL)
     AND    (meta.dictionary.name = $4::varchar OR $4::varchar IS NULL);`

	ServerLink = `
SELECT      asset.server_linking.serverID_A as serverID,
            asset.server_linking.dictionaryID_A as dictionaryID
FROM        asset.server_linking
WHERE       asset.server_linking.serverID_B = $1:uuid
UNION
SELECT      asset.server_linking.serverID_B as serverID,
            asset.server_linking.dictionaryID_B as dictionaryID
FROM        asset.server_linking
WHERE       asset.server_linking.serverID_A = $1:uuid;`

	ServerList = `
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

	ServerParent = `
SELECT      asset.runtime_environment.rteID,
            asset.runtime_environment.dictionaryID,
            meta.dictionary.name,
            asset.runtime_environment_unique_attribute_values.value
FROM        asset.server
    JOIN    asset.server_parent
      ON    asset.server.serverID = asset.server_parent.serverID
    JOIN    asset.runtime_environment
      ON    asset.server_parent.parentRuntimeID = asset.runtime_environment.rteID
    JOIN    meta.dictionary
      ON    asset.runtime_environment.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.runtime_environment_unique_attribute_values
      ON    asset.runtime_environment_unique_attribute_values.rteID = asset.runtime_environment.rteID
     AND    asset.runtime_environment_unique_attribute_values.dictionaryID = asset.runtime_environment.dictionaryID
     AND    asset.runtime_environment_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
WHERE       asset.server.serverID = $1::uuid
     AND    now()::timestamptz(3) <@ asset.server_parent.validity
     AND    now()::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity
     AND    meta.unique_attribute.attribute IN ('name');`
)

func init() {
	m[ServerAttribute] = `ServerAttribute`
	m[ServerFind] = `ServerFind`
	m[ServerLink] = `ServerLink`
	m[ServerList] = `ServerList`
	m[ServerParent] = `ServerParent`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
