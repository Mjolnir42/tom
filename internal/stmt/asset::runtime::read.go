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
SELECT            meta.dictionary.name AS dictionaryName,
                  asset.runtime_environment_unique_attribute_values.value AS runtimeName,
                  inventory.user.uid AS createdBy,
                  asset.runtime_environment_unique_attribute_values.createdAt
FROM              meta.dictionary
JOIN              meta.unique_attribute
  ON              meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
JOIN              asset.runtime_environment
  ON              meta.dictionary.dictionaryID = asset.runtime_environment.dictionaryID
JOIN              asset.runtime_environment_unique_attribute_values
    ON            meta.dictionary.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
    AND           asset.runtime_environment.rteID = asset.runtime_environment_unique_attribute_values.rteID
    AND           meta.unique_attribute.attributeID = asset.runtime_environment_unique_attribute_values.attributeID
JOIN              inventory.user
  ON              asset.runtime_environment_unique_attribute_values.createdBy = inventory.user.userID
WHERE             (meta.dictionary.name = $1::text OR $1::text IS NULL)
  AND             meta.unique_attribute.attribute = 'name'::text
  AND             now()::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	RuntimeTxShow = `
SELECT            asset.runtime_environment.rteID,
                  asset.runtime_environment.dictionaryID,
                  asset.runtime_environment.createdAt,
                  creator.uid AS createdBy,
                  lower(asset.runtime_environment_unique_attribute_values.validity) AS validSince,
                  upper(asset.runtime_environment_unique_attribute_values.validity) AS validUntil,
                  asset.runtime_environment_unique_attribute_values.createdAt AS namedAt,
                  namegiver.uid AS namedBy
FROM              meta.dictionary
    JOIN          asset.runtime_environment
        ON        meta.dictionary.dictionaryID = asset.runtime_environment.dictionaryID
    JOIN          inventory.user AS creator
        ON        asset.runtime_environment.createdBy = creator.userID
    JOIN          meta.unique_attribute
        ON        meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN          asset.runtime_environment_unique_attribute_values
        ON        meta.dictionary.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
        AND       asset.runtime_environment.rteID = asset.runtime_environment_unique_attribute_values.rteID
        AND       meta.unique_attribute.attributeID = asset.runtime_environment_unique_attribute_values.attributeID
    JOIN          inventory.user AS namegiver
        ON        asset.runtime_environment_unique_attribute_values.createdBy = namegiver.userID
WHERE             meta.dictionary.name = $1::text
     AND          meta.unique_attribute.attribute = 'name'::text
     AND          asset.runtime_environment_unique_attribute_values.value = $2::text
     AND          $3::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	RuntimeListLinked = `
WITH sel_rte AS ( SELECT linkedViaA.rteID_B AS linkedRteID,
                         linkedViaA.dictionaryID_B AS linkedDictID
                  FROM   asset.runtime_environment
                  JOIN   asset.runtime_environment_linking AS linkedViaA
                    ON   asset.runtime_environment.rteID = linkedViaA.rteID_A
                  WHERE  asset.runtime_environment.rteID = $1::uuid
                    AND  asset.runtime_environment.dictionaryID = $2::uuid
                  UNION
                  SELECT linkedViaB.rteID_A AS linkedRteID,
                         linkedViaB.dictionaryID_A AS linkedDictID
                  FROM   asset.runtime_environment
                  JOIN   asset.runtime_environment_linking AS linkedViaB
                    ON   asset.runtime_environment.rteID = linkedViaB.rteID_B
                  WHERE  asset.runtime_environment.rteID = $1::uuid
                    AND  asset.runtime_environment.dictionaryID = $2::uuid)
SELECT            sel_rte.linkedRteID AS rteID,
                  sel_rte.linkedDictID AS dictionaryID,
                  asset.runtime_environment_unique_attribute_values.value AS name,
                  meta.dictionary.name AS namespace
FROM              sel_rte
JOIN              asset.runtime_environment
  ON              sel_rte.linkedRteID
   =              asset.runtime_environment.rteID
 AND              sel_rte.linkedDictID
   =              asset.runtime_environment.dictionaryID
JOIN              meta.unique_attribute
  ON              asset.runtime_environment.dictionaryID
   =              meta.unique_attribute.dictionaryID
JOIN              asset.runtime_environment_unique_attribute_values
  ON              sel_rte.linkedRteID
   =              asset.runtime_environment_unique_attribute_values.rteID
 AND              sel_rte.linkedDictID
   =              asset.runtime_environment_unique_attribute_values.dictionaryID
 AND              meta.unique_attribute.attributeID
   =              asset.runtime_environment_unique_attribute_values.attributeID
JOIN              meta.dictionary
  ON              sel_rte.linkedDictID = meta.dictionary.dictionaryID
WHERE             meta.unique_attribute.attribute = 'name'::text
  AND             $3::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	RuntimeTxShowProperties = `
SELECT      meta.unique_attribute.attribute AS attribute,
            asset.runtime_environment_unique_attribute_values.value AS value,
            lower(asset.runtime_environment_unique_attribute_values.validity) AS validSince,
            upper(asset.runtime_environment_unique_attribute_values.validity) AS validUntil,
            asset.runtime_environment_unique_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.runtime_environment_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.runtime_environment_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.runtime_environment_unique_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.runtime_environment_unique_attribute_values.rteID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            asset.runtime_environment_standard_attribute_values.value AS value,
            lower(asset.runtime_environment_standard_attribute_values.validity) AS validSince,
            upper(asset.runtime_environment_standard_attribute_values.validity) AS validUntil,
            asset.runtime_environment_standard_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.runtime_environment_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = asset.runtime_environment_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = asset.runtime_environment_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.runtime_environment_standard_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.runtime_environment_standard_attribute_values.rteID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.runtime_environment_standard_attribute_values.validity;`

	RuntimeParent = `
SELECT      'runtime'::text                                         AS entity,
            asset.runtime_environment.rteID                         AS objID,
            asset.runtime_environment.dictionaryID                  AS dictID,
            meta.dictionary.name                                    AS dictName,
            asset.runtime_environment_unique_attribute_values.value AS objName
FROM        asset.runtime_environment
    JOIN    asset.runtime_environment_parent
      ON    asset.runtime_environment.rteID = asset.runtime_environment_parent.rteID
    JOIN    asset.runtime_environment AS parent
      ON    asset.runtime_environment_parent.parentRuntimeID = parent.rteID
    JOIN    meta.dictionary
      ON    parent.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.runtime_environment_unique_attribute_values
      ON    asset.runtime_environment_unique_attribute_values.rteID = parent.rteID
     AND    asset.runtime_environment_unique_attribute_values.dictionaryID = parent.dictionaryID
     AND    asset.runtime_environment_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
            -- runtime environment for which the parent is searched
WHERE       asset.runtime_environment.rteID = $1::uuid
            -- parent relationship is currently valid
     AND    $2::timestamptz(3) <@ asset.runtime_environment_parent.validity
            -- registered parent is still valid, based on the validity of the parent's name
     AND    $2::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity
     AND    meta.unique_attribute.attribute IN ('name')
UNION
SELECT      'server'::text                             AS entity,
            asset.server.serverID                      AS objID,
            asset.server.dictionaryID                  AS dictID,
            meta.dictionary.name                       AS dictName,
            asset.server_unique_attribute_values.value AS objName
FROM        asset.runtime_environment
    JOIN    asset.runtime_environment_parent
      ON    asset.runtime_environment.rteID = asset.runtime_environment_parent.rteID
    JOIN    asset.server
      ON    asset.runtime_environment_parent.parentServerID = asset.server.serverID
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    asset.server_unique_attribute_values.serverID = asset.server.serverID
     AND    asset.server_unique_attribute_values.dictionaryID = asset.server.dictionaryID
     AND    asset.server_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
            -- runtime environment for which the parent is searched
WHERE       asset.runtime_environment.rteID = $1::uuid
            -- parent relationship is currently valid
     AND    $2::timestamptz(3) <@ asset.runtime_environment_parent.validity
            -- registered parent is still valid, based on the validity of the parent's name
     AND    $2::timestamptz(3) <@ asset.server_unique_attribute_values.validity
     AND    meta.unique_attribute.attribute IN ('name')
UNION
SELECT      'orchestration'::text                                         AS entity,
            asset.orchestration_environment.orchID                        AS objID,
            asset.orchestration_environment.dictionaryID                  AS dictID,
            meta.dictionary.name                                          AS dictName,
            asset.orchestration_environment_unique_attribute_values.value AS objName
FROM        asset.runtime_environment
    JOIN    asset.runtime_environment_parent
      ON    asset.runtime_environment.rteID = asset.runtime_environment_parent.rteID
    JOIN    asset.orchestration_environment
      ON    asset.runtime_environment_parent.parentOrchestrationID = asset.orchestration_environment.orchID
    JOIN    meta.dictionary
      ON    asset.orchestration_environment.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.orchestration_environment_unique_attribute_values
      ON    asset.orchestration_environment_unique_attribute_values.orchID = asset.orchestration_environment.orchID
     AND    asset.orchestration_environment_unique_attribute_values.dictionaryID = asset.orchestration_environment.dictionaryID
     AND    asset.orchestration_environment_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
            -- runtime environment for which the parent is searched
WHERE       asset.runtime_environment.rteID = $1::uuid
            -- parent relationship is currently valid
     AND    $2::timestamptz(3) <@ asset.runtime_environment_parent.validity
            -- registered parent is still valid, based on the validity of the parent's name
     AND    $2::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity
     AND    meta.unique_attribute.attribute IN ('name');`

	RuntimeTxShowChildren = `
SELECT      'server'::text AS childEntity,
            asset.server_unique_attribute_values.value AS childName,
            meta.dictionary.name AS childDictName
FROM        asset.server_parent
    JOIN    asset.server
      ON    asset.server_parent.serverID = asset.server.serverID
    JOIN    meta.dictionary
      ON    asset.server.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.server_unique_attribute_values
      ON    asset.server.serverID = asset.server_unique_attribute_values.serverID
     AND    meta.unique_attribute.dictionaryID = asset.server_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.server_unique_attribute_values.attributeID
WHERE       asset.server_parent.parentRuntimeID = $1::uuid
     AND    meta.unique_attribute.attribute = 'name'::text
     AND    $2::timestamptz(3) <@ asset.server_parent.validity
     AND    $2::timestamptz(3) <@ asset.server_unique_attribute_values.validity
UNION
SELECT      'orchestration'::text AS childEntity,
            asset.orchestration_environment_unique_attribute_values.value AS childName,
            meta.dictionary.name AS childDictName
FROM        asset.orchestration_environment_mapping
    JOIN    asset.orchestration_environment
      ON    asset.orchestration_environment_mapping.orchID = asset.orchestration_environment.orchID
    JOIN    meta.dictionary
      ON    asset.orchestration_environment.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.orchestration_environment_unique_attribute_values
      ON    asset.orchestration_environment.orchID = asset.orchestration_environment_unique_attribute_values.orchID
     AND    meta.unique_attribute.dictionaryID = asset.orchestration_environment_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.orchestration_environment_unique_attribute_values.attributeID
WHERE       asset.orchestration_environment_mapping.parentRuntimeID = $1::uuid
     AND    meta.unique_attribute.attribute = 'name'::text
     AND    $2::timestamptz(3) <@ asset.orchestration_environment_mapping.validity
     AND    $2::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity
UNION
SELECT      'runtime'::text AS childEntity,
            asset.runtime_environment_unique_attribute_values.value AS childName,
            meta.dictionary.name AS childDictName
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
WHERE       asset.runtime_environment_parent.parentRuntimeID = $1::uuid
     AND    meta.unique_attribute.attribute = 'name'::text
     AND    $2::timestamptz(3) <@ asset.runtime_environment_parent.validity
     AND    $2::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	RuntimeResolveServer = `
SELECT      asset.server_unique_attribute_values.value,
            meta.dictionary.name,
            resolution.serverType
FROM        view.resolveRuntimeToServer($1::uuid) AS resolution
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

	RuntimeResolvePhysical = `
SELECT      asset.server_unique_attribute_values.value,
            meta.dictionary.name,
            resolution.serverType
FROM        view.resolveRuntimeToPhysical($1::uuid) AS resolution
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

	RuntimeTxSelectResource = `
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
      uri AS ( SELECT meta.dictionary_standard_attribute_values.value AS uri
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
SELECT                replace(uri.uri, '{{LOOKUP}}', asset.runtime_environment_unique_attribute_values.value) AS resource
FROM                  asset.runtime_environment
JOIN                  dict
  ON                  asset.runtime_environment.dictionaryID = dict.dictionaryID
JOIN                  meta.unique_attribute
  ON                  asset.runtime_environment.dictionaryID = meta.unique_attribute.dictionaryID
JOIN                  asset.runtime_environment_unique_attribute_values
  ON                  asset.runtime_environment.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
 AND                  meta.unique_attribute.attributeID = asset.runtime_environment_unique_attribute_values.attributeID
 AND                  asset.runtime_environment.rteID = asset.runtime_environment_unique_attribute_values.rteID
JOIN                  look
  ON                  meta.unique_attribute.attribute = look.key
CROSS JOIN            uri
WHERE                 asset.runtime_environment.rteID = $2::uuid
  AND                 $3::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`
)

func init() {
	m[RuntimeListLinked] = `RuntimeListLinked`
	m[RuntimeList] = `RuntimeList`
	m[RuntimeParent] = `RuntimeParent`
	m[RuntimeResolvePhysical] = `RuntimeResolvePhysical`
	m[RuntimeResolveServer] = `RuntimeResolveServer`
	m[RuntimeTxSelectResource] = `RuntimeTxSelectResource`
	m[RuntimeTxShowChildren] = `RuntimeTxShowChildren`
	m[RuntimeTxShowProperties] = `RuntimeTxShowProperties`
	m[RuntimeTxShow] = `RuntimeTxShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
