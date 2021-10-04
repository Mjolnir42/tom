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
)

func init() {
	m[RuntimeList] = `RuntimeList`
	m[RuntimeListLinked] = `RuntimeListLinked`
	m[RuntimeTxShow] = `RuntimeTxShow`
	m[RuntimeTxShowProperties] = `RuntimeTxShowProperties`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
