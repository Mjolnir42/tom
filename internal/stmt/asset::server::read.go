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

	ServerList = `
SELECT      asset.server.serverID,
            meta.dictionary.name,
            meta.standard_attribute.attribute,
            asset.server_standard_attribute_values.value,
            inventory.user.uid AS createdBy,
            asset.server_standard_attribute_values.createdAt
FROM        asset.server
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.server_standard_attribute_values
      ON    asset.server_standard_attribute_values.serverID = asset.server.serverID
     AND    asset.server_standard_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_standard_attribute_values.attributeID = meta.standard_attribute.attributeID
    JOIN    inventory.user
      ON    asset.server_standard_attribute_values.createdBy = inventory.user.userID
   WHERE    now()::timestamptz(3) <@ asset.server_standard_attribute_values.validity
     AND    (meta.dictionary.name = $1::text OR $1::text IS NULL)
     AND    meta.standard_attribute.attribute IN ('type')
UNION
SELECT      asset.server.serverID,
            meta.dictionary.name,
            meta.unique_attribute.attribute,
            asset.server_unique_attribute_values.value,
            inventory.user.uid AS createdBy,
            asset.server_unique_attribute_values.createdAt
FROM        asset.server
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    asset.server_unique_attribute_values.serverID = asset.server.serverID
     AND    asset.server_unique_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
    JOIN    inventory.user
      ON    asset.server_unique_attribute_values.createdBy = inventory.user.userID
WHERE       now()::timestamptz(3) <@ asset.server_unique_attribute_values.validity
     AND    (meta.dictionary.name = $1::text OR $1::text IS NULL)
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

	ServerTxShow = `
SELECT            asset.server.serverID,
                  asset.server.dictionaryID,
                  asset.server.createdAt,
                  creator.uid AS createdBy,
                  lower(asset.server_unique_attribute_values.validity) AS validSince,
                  upper(asset.server_unique_attribute_values.validity) AS validUntil,
                  asset.server_unique_attribute_values.createdAt AS namedAt,
                  namegiver.uid AS namedBy
FROM              meta.dictionary
    JOIN          asset.server
        ON        meta.dictionary.dictionaryID = asset.server.dictionaryID
    JOIN          inventory.user AS creator
        ON        asset.server.createdBy = creator.userID
    JOIN          meta.unique_attribute
        ON        meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN          asset.server_unique_attribute_values
        ON        meta.dictionary.dictionaryID = asset.server_unique_attribute_values.dictionaryID
        AND       asset.server.serverID = asset.server_unique_attribute_values.serverID
        AND       meta.unique_attribute.attributeID = asset.server_unique_attribute_values.attributeID
    JOIN          inventory.user AS namegiver
        ON        asset.server_unique_attribute_values.createdBy = namegiver.userID
WHERE             meta.dictionary.name = $1::text
     AND          meta.unique_attribute.attribute = 'name'::text
     AND          asset.server_unique_attribute_values.value = $2::text
     AND          $3::timestamptz(3) <@ asset.server_unique_attribute_values.validity;`

	ServerListLinked = `
WITH sel_cte AS ( SELECT linkedViaA.serverID_B AS linkedServerID,
                         linkedViaA.dictionaryID_B AS linkedDictID
                  FROM   asset.server
                  JOIN   asset.server_linking AS linkedViaA
                    ON   asset.server.serverID = linkedViaA.serverID_A
                  WHERE  asset.server.serverID = $1::uuid
                    AND  asset.server.dictionaryID = $2::uuid
                  UNION
                  SELECT linkedViaB.serverID_A AS linkedServerID,
                         linkedViaB.dictionaryID_A AS linkedDictID
                  FROM   asset.server
                  JOIN   asset.server_linking AS linkedViaB
                    ON   asset.server.serverID = linkedViaB.serverID_B
                  WHERE  asset.server.serverID = $1::uuid
                    AND  asset.server.dictionaryID = $2::uuid)
SELECT            sel_cte.linkedServerID AS serverID,
                  sel_cte.linkedDictID AS dictionaryID,
                  asset.server_unique_attribute_values.value AS name,
                  meta.dictionary.name AS namespace
FROM              sel_cte
JOIN              asset.server
  ON              sel_cte.linkedServerID
   =              asset.server.serverID
 AND              sel_cte.linkedDictID
   =              asset.server.dictionaryID
JOIN              meta.unique_attribute
  ON              asset.server.dictionaryID
   =              meta.unique_attribute.dictionaryID
JOIN              asset.server_unique_attribute_values
  ON              sel_cte.linkedServerID
   =              asset.server_unique_attribute_values.serverID
 AND              sel_cte.linkedDictID
   =              asset.server_unique_attribute_values.dictionaryID
 AND              meta.unique_attribute.attributeID
   =              asset.server_unique_attribute_values.attributeID
JOIN              meta.dictionary
  ON              sel_cte.linkedDictID = meta.dictionary.dictionaryID
WHERE             meta.unique_attribute.attribute = 'name'::text
  AND             $3::timestamptz(3) <@ asset.server_unique_attribute_values.validity;`

	ServerTxShowProperties = `
SELECT      meta.unique_attribute.attribute AS attribute,
            asset.server_unique_attribute_values.value AS value,
            lower(asset.server_unique_attribute_values.validity) AS validSince,
            upper(asset.server_unique_attribute_values.validity) AS validUntil,
            asset.server_unique_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = asset.server_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.server_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.server_unique_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.server_unique_attribute_values.serverID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.server_unique_attribute_values.validity
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            asset.server_standard_attribute_values.value AS value,
            lower(asset.server_standard_attribute_values.validity) AS validSince,
            upper(asset.server_standard_attribute_values.validity) AS validUntil,
            asset.server_standard_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.server_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = asset.server_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = asset.server_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.server_standard_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.server_standard_attribute_values.serverID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.server_standard_attribute_values.validity;`

	ServerTxShowChildren = `
SELECT      asset.runtime_environment_unique_attribute_values.value,
            meta.dictionary.name
FROM        asset.runtime_environment_parent
    JOIN    asset.runtime_environment
      ON    asset.runtime_environment_parent.rteID = asset.runtime_environment.rteID
    JOIN    meta.dictionary
      ON    asset.runtime_environment.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.runtime_environment_unique_attribute_values
      ON    asset.runtime_environment.rteID = asset.runtime_environment_unique_attribute_values.rteID
		 AND    meta.unique_attribute.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.runtime_environment_unique_attribute_values.attributeID
WHERE       asset.runtime_environment_parent.parentServerID = $1::uuid
     AND    meta.unique_attribute.attribute = 'name'::text
     AND    $2::timestamptz(3) <@ asset.runtime_environment_parent.validity
     AND    $2::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	ServerResolveServer = `
SELECT      asset.server_unique_attribute_values.value,
            meta.dictionary.name,
            resolution.serverType
FROM        view.resolveServerToServer($1::uuid) AS resolution
    JOIN    asset.server
      ON    resolution.serverID = asset.server.serverID
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    asset.server_unique_attribute_values.serverID = asset.server.serverID
     AND    asset.server_unique_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
WHERE       meta.unique_attribute.attribute IN ('name');`

	ServerResolvePhysical = `
SELECT      asset.server_unique_attribute_values.value,
            meta.dictionary.name,
            resolution.serverType
FROM        view.resolveServerToPhysical($1::uuid) AS resolution
    JOIN    asset.server
      ON    resolution.serverID = asset.server.serverID
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    asset.server_unique_attribute_values.serverID = asset.server.serverID
     AND    asset.server_unique_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
WHERE       meta.unique_attribute.attribute IN ('name');`

	ServerTxSelectResource = `
WITH dict AS ( SELECT meta.dictionary.dictionaryID
               FROM   meta.dictionary
               JOIN   meta.standard_attribute
                 ON   meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
               JOIN   meta.dictionary_standard_attribute_values
                 ON   meta.dictionary.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
                AND   meta.standard_attribute.attributeID = meta.dictionary_standard_attribute_values.attributeID
               WHERE  meta.dictionary.name = $1::text
                 AND  meta.standard_attribute.attribute = 'dict_type'
                 AND  meta.dictionary_standard_attribute_values.value = 'referential'
                 AND  $3::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity),
     look AS ( SELECT meta.dictionary_standard_attribute_values.value AS key
               FROM   meta.dictionary
               JOIN   dict
                 ON   dict.dictionaryID = meta.dictionary.dictionaryID
               JOIN   meta.standard_attribute
                 ON   meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
               JOIN   meta.dictionary_standard_attribute_values
                 ON   meta.dictionary.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
                AND   meta.standard_attribute.attributeID = meta.dictionary_standard_attribute_values.attributeID
               WHERE  meta.standard_attribute.attribute = 'dict_lookup'
                 AND  $3::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity),
      uri AS ( SELECT meta.dictionary_standard_attribute_values.value AS uri,
                      meta.dictionary.dictionaryID
               FROM   meta.dictionary
               JOIN   dict
                 ON   dict.dictionaryID = meta.dictionary.dictionaryID
               JOIN   meta.standard_attribute
                 ON   meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
               JOIN   meta.dictionary_standard_attribute_values
                 ON   meta.dictionary.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
                AND   meta.standard_attribute.attributeID = meta.dictionary_standard_attribute_values.attributeID
               WHERE  meta.standard_attribute.attribute = 'dict_uri'
                 AND  $3::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity )
SELECT                replace(uri.uri, '{{LOOKUP}}', asset.server_unique_attribute_values.value) AS resource
FROM                  asset.server
JOIN                  dict
  ON                  asset.server.dictionaryID = dict.dictionaryID
JOIN                  meta.unique_attribute
  ON                  asset.server.dictionaryID = meta.unique_attribute.dictionaryID
JOIN                  asset.server_unique_attribute_values
  ON                  asset.server.dictionaryID = asset.server_unique_attribute_values.dictionaryID
 AND                  meta.unique_attribute.attributeID = asset.server_unique_attribute_values.attributeID
 AND                  asset.server.serverID = asset.server_unique_attribute_values.serverID
JOIN                  look
  ON                  meta.unique_attribute.attribute = look.key
JOIN                  uri
  ON                  asset.server.dictionaryID = uri.dictionaryID
WHERE                 asset.server.serverID = $2::uuid
  AND                 $3::timestamptz(3) <@ asset.server_unique_attribute_values.validity;`
)

func init() {
	m[ServerAttribute] = `ServerAttribute`
	m[ServerFind] = `ServerFind`
	m[ServerListLinked] = `ServerListLinked`
	m[ServerList] = `ServerList`
	m[ServerParent] = `ServerParent`
	m[ServerResolvePhysical] = `ServerResolvePhysical`
	m[ServerResolveServer] = `ServerResolveServer`
	m[ServerTxSelectResource] = `ServerTxSelectResource`
	m[ServerTxShowChildren] = `ServerTxShowChildren`
	m[ServerTxShowProperties] = `ServerTxShowProperties`
	m[ServerTxShow] = `ServerTxShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
