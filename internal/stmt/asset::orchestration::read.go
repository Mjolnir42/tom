/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	OrchestrationReadStatements = ``

	OrchestrationListLinked = `
WITH sel_cte AS ( SELECT linkedViaA.orchID_B AS linkedOrchestrationID,
                         linkedViaA.dictionaryID_B AS linkedDictID
                  FROM   asset.orchestration_environment
                  JOIN   asset.orchestration_environment_linking AS linkedViaA
                    ON   asset.orchestration_environment.orchID = linkedViaA.orchID_A
                  WHERE  asset.orchestration_environment.orchID = $1::uuid
                    AND  asset.orchestration_environment.dictionaryID = $2::uuid
                  UNION
                  SELECT linkedViaB.orchID_A AS linkedOrchestrationID,
                         linkedViaB.dictionaryID_A AS linkedDictID
                  FROM   asset.orchestration_environment
                  JOIN   asset.orchestration_environment_linking AS linkedViaB
                    ON   asset.orchestration_environment.orchID = linkedViaB.orchID_B
                  WHERE  asset.orchestration_environment.orchID = $1::uuid
                    AND  asset.orchestration_environment.dictionaryID = $2::uuid)
SELECT            sel_cte.linkedOrchestrationID AS orchID,
                  sel_cte.linkedDictID AS dictionaryID,
                  asset.orchestration_environment_unique_attribute_values.value AS name,
                  meta.dictionary.name AS namespace
FROM              sel_cte
JOIN              asset.orchestration_environment
  ON              sel_cte.linkedOrchestrationID
   =              asset.orchestration_environment.orchID
 AND              sel_cte.linkedDictID
   =              asset.orchestration_environment.dictionaryID
JOIN              meta.unique_attribute
  ON              asset.orchestration_environment.dictionaryID
   =              meta.unique_attribute.dictionaryID
JOIN              asset.orchestration_environment_unique_attribute_values
  ON              sel_cte.linkedOrchestrationID
   =              asset.orchestration_environment_unique_attribute_values.orchID
 AND              sel_cte.linkedDictID
   =              asset.orchestration_environment_unique_attribute_values.dictionaryID
 AND              meta.unique_attribute.attributeID
   =              asset.orchestration_environment_unique_attribute_values.attributeID
JOIN              meta.dictionary
  ON              sel_cte.linkedDictID = meta.dictionary.dictionaryID
WHERE             meta.unique_attribute.attribute = 'name'::text
  AND             $3::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity;`

	OrchestrationList = `
SELECT      asset.orchestration_environment.orchID,
            meta.dictionary.name,
            meta.standard_attribute.attribute,
            asset.orchestration_environment_standard_attribute_values.value,
            inventory.user.uid AS createdBy,
            asset.orchestration_environment_standard_attribute_values.createdAt
FROM        asset.orchestration_environment
    JOIN    meta.dictionary
      ON    asset.orchestration_environment.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.orchestration_environment_standard_attribute_values
      ON    asset.orchestration_environment_standard_attribute_values.orchID = asset.orchestration_environment.orchID
     AND    asset.orchestration_environment_standard_attribute_values.dictionaryID = asset.orchestration_environment.dictionaryID
     AND    asset.orchestration_environment_standard_attribute_values.attributeID = meta.standard_attribute.attributeID
    JOIN    inventory.user
      ON    asset.orchestration_environment_standard_attribute_values.createdBy = inventory.user.userID
   WHERE    now()::timestamptz(3) <@ asset.orchestration_environment_standard_attribute_values.validity
	   AND    (meta.dictionary.name = $1::text OR $1::text IS NULL)
     AND    meta.standard_attribute.attribute IN ('type')
UNION
SELECT      asset.orchestration_environment.orchID,
            meta.dictionary.name,
            meta.unique_attribute.attribute,
            asset.orchestration_environment_unique_attribute_values.value,
            inventory.user.uid AS createdBy,
            asset.orchestration_environment_unique_attribute_values.createdAt
FROM        asset.orchestration_environment
    JOIN    meta.dictionary
      ON    asset.orchestration_environment.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.orchestration_environment_unique_attribute_values
      ON    asset.orchestration_environment_unique_attribute_values.orchID = asset.orchestration_environment.orchID
     AND    asset.orchestration_environment_unique_attribute_values.dictionaryID = asset.orchestration_environment.dictionaryID
     AND    asset.orchestration_environment_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
    JOIN    inventory.user
      ON    asset.orchestration_environment_unique_attribute_values.createdBy = inventory.user.userID
WHERE       now()::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity
     AND    (meta.dictionary.name = $1::text OR $1::text IS NULL)
     AND    meta.unique_attribute.attribute IN ('name');`

	OrchestrationTxParent = `
SELECT      parent.rteID,
            parent.dictionaryID,
            meta.dictionary.name,
            asset.runtime_environment_unique_attribute_values.value
FROM        asset.orchestration_environment
    JOIN    asset.orchestration_environment_mapping
      ON    asset.orchestration_environment.orchID = asset.orchestration_environment_mapping.orchID
    JOIN    asset.runtime_environment AS parent
      ON    asset.orchestration_environment_mapping.parentRuntimeID = parent.rteID
    JOIN    meta.dictionary
      ON    parent.dictionaryID = meta.dictionary.dictionaryID
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.runtime_environment_unique_attribute_values
      ON    asset.runtime_environment_unique_attribute_values.rteID = parent.rteID
     AND    asset.runtime_environment_unique_attribute_values.dictionaryID = parent.dictionaryID
     AND    asset.runtime_environment_unique_attribute_values.attributeID = meta.unique_attribute.attributeID
            -- orchestration environment for which the parent is searched
WHERE       asset.orchestration_environment.orchID = $1::uuid
            -- parent relationship is currently valid
     AND    $2::timestamptz(3) <@ asset.orchestration_environment_mapping.validity
            -- registered parent is still valid, based on the validity of the parent's name
     AND    $2::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity
     AND    meta.unique_attribute.attribute IN ('name');`

	OrchestrationTxShowChildren = `
SELECT      asset.runtime_environment_unique_attribute_values.value AS childName,
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
WHERE       asset.runtime_environment_parent.parentOrchestrationID = $1::uuid
     AND    meta.unique_attribute.attribute = 'name'::text
     AND    $2::timestamptz(3) <@ asset.runtime_environment_parent.validity
     AND    $2::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	OrchestrationTxShowProperties = `
SELECT      meta.unique_attribute.attribute AS attribute,
            asset.orchestration_environment_unique_attribute_values.value AS value,
            lower(asset.orchestration_environment_unique_attribute_values.validity) AS validSince,
            upper(asset.orchestration_environment_unique_attribute_values.validity) AS validUntil,
            asset.orchestration_environment_unique_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.orchestration_environment_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = asset.orchestration_environment_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.orchestration_environment_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.orchestration_environment_unique_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.orchestration_environment_unique_attribute_values.orchID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            asset.orchestration_environment_standard_attribute_values.value AS value,
            lower(asset.orchestration_environment_standard_attribute_values.validity) AS validSince,
            upper(asset.orchestration_environment_standard_attribute_values.validity) AS validUntil,
            asset.orchestration_environment_standard_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.orchestration_environment_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = asset.orchestration_environment_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = asset.orchestration_environment_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.orchestration_environment_standard_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.orchestration_environment_standard_attribute_values.orchID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.orchestration_environment_standard_attribute_values.validity;`

	OrchestrationTxShow = `
SELECT            asset.orchestration_environment.orchID,
                  asset.orchestration_environment.dictionaryID,
                  asset.orchestration_environment.createdAt,
                  creator.uid AS createdBy,
                  lower(asset.orchestration_environment_unique_attribute_values.validity) AS validSince,
                  upper(asset.orchestration_environment_unique_attribute_values.validity) AS validUntil,
                  asset.orchestration_environment_unique_attribute_values.createdAt AS namedAt,
                  namegiver.uid AS namedBy
FROM              meta.dictionary
    JOIN          asset.orchestration_environment
        ON        meta.dictionary.dictionaryID = asset.orchestration_environment.dictionaryID
    JOIN          inventory.user AS creator
        ON        asset.orchestration_environment.createdBy = creator.userID
    JOIN          meta.unique_attribute
        ON        meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN          asset.orchestration_environment_unique_attribute_values
        ON        meta.dictionary.dictionaryID = asset.orchestration_environment_unique_attribute_values.dictionaryID
        AND       asset.orchestration_environment.orchID = asset.orchestration_environment_unique_attribute_values.orchID
        AND       meta.unique_attribute.attributeID = asset.orchestration_environment_unique_attribute_values.attributeID
    JOIN          inventory.user AS namegiver
        ON        asset.orchestration_environment_unique_attribute_values.createdBy = namegiver.userID
WHERE             meta.dictionary.name = $1::text
     AND          meta.unique_attribute.attribute = 'name'::text
     AND          asset.orchestration_environment_unique_attribute_values.value = $2::text
     AND          $3::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity;`
)

func init() {
	m[OrchestrationListLinked] = `OrchestrationListLinked`
	m[OrchestrationList] = `OrchestrationList`
	m[OrchestrationTxParent] = `OrchestrationTxParent`
	m[OrchestrationTxShowChildren] = `OrchestrationTxShowChildren`
	m[OrchestrationTxShowProperties] = `OrchestrationTxShowProperties`
	m[OrchestrationTxShow] = `OrchestrationTxShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
