/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	FlowReadStatements = ``

	FlowList = `
SELECT            bulk.flow.flowID,
                  meta.dictionary.name AS dictionaryName,
                  meta.unique_attribute.attribute,
                  bulk.flow_unique_attribute_values.value,
                  inventory.user.uid AS createdBy,
                  bulk.flow_unique_attribute_values.createdAt
FROM              meta.dictionary
JOIN              meta.unique_attribute
  ON              meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
JOIN              bulk.flow
  ON              meta.dictionary.dictionaryID = bulk.flow.dictionaryID
JOIN              bulk.flow_unique_attribute_values
    ON            meta.dictionary.dictionaryID = bulk.flow_unique_attribute_values.dictionaryID
    AND           bulk.flow.flowID = bulk.flow_unique_attribute_values.flowID
    AND           meta.unique_attribute.attributeID = bulk.flow_unique_attribute_values.attributeID
JOIN              inventory.user
  ON              bulk.flow_unique_attribute_values.createdBy = inventory.user.userID
WHERE             (meta.dictionary.name = $1::text OR $1::text IS NULL)
  AND             meta.unique_attribute.attribute = 'name'::text
  AND             now()::timestamptz(3) <@ bulk.flow_unique_attribute_values.validity
UNION
SELECT            bulk.flow.flowID,
                  meta.dictionary.name AS dictionaryName,
                  meta.standard_attribute.attribute,
                  bulk.flow_standard_attribute_values.value,
                  inventory.user.uid AS createdBy,
                  bulk.flow_standard_attribute_values.createdAt
FROM              meta.dictionary
JOIN              meta.standard_attribute
  ON              meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
JOIN              bulk.flow
  ON              meta.dictionary.dictionaryID = bulk.flow.dictionaryID
JOIN              bulk.flow_standard_attribute_values
    ON            meta.dictionary.dictionaryID = bulk.flow_standard_attribute_values.dictionaryID
    AND           bulk.flow.flowID = bulk.flow_standard_attribute_values.flowID
    AND           meta.standard_attribute.attributeID = bulk.flow_standard_attribute_values.attributeID
JOIN              inventory.user
  ON              bulk.flow_standard_attribute_values.createdBy = inventory.user.userID
WHERE             (meta.dictionary.name = $1::text OR $1::text IS NULL)
  AND             meta.standard_attribute.attribute = 'type'::text
  AND             now()::timestamptz(3) <@ bulk.flow_standard_attribute_values.validity;`

	FlowTxShow = `
SELECT            bulk.flow.flowID,
                  bulk.flow.dictionaryID,
                  bulk.flow.createdAt,
                  creator.uid AS createdBy,
                  lower(bulk.flow_unique_attribute_values.validity) AS validSince,
                  upper(bulk.flow_unique_attribute_values.validity) AS validUntil,
                  bulk.flow_unique_attribute_values.createdAt AS namedAt,
                  namegiver.uid AS namedBy
FROM              meta.dictionary
    JOIN          bulk.flow
        ON        meta.dictionary.dictionaryID = bulk.flow.dictionaryID
    JOIN          inventory.user AS creator
        ON        bulk.flow.createdBy = creator.userID
    JOIN          meta.unique_attribute
        ON        meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN          bulk.flow_unique_attribute_values
        ON        meta.dictionary.dictionaryID = bulk.flow_unique_attribute_values.dictionaryID
        AND       bulk.flow.flowID = bulk.flow_unique_attribute_values.flowID
        AND       meta.unique_attribute.attributeID = bulk.flow_unique_attribute_values.attributeID
    JOIN          inventory.user AS namegiver
        ON        bulk.flow_unique_attribute_values.createdBy = namegiver.userID
WHERE             meta.dictionary.name = $1::text
     AND          meta.unique_attribute.attribute = 'name'::text
     AND          bulk.flow_unique_attribute_values.value = $2::text
     AND          $3::timestamptz(3) <@ bulk.flow_unique_attribute_values.validity;`

	FlowTxShowProperties = `
SELECT      meta.unique_attribute.attribute AS attribute,
            bulk.flow_unique_attribute_values.value AS value,
            lower(bulk.flow_unique_attribute_values.validity) AS validSince,
            upper(bulk.flow_unique_attribute_values.validity) AS validUntil,
            bulk.flow_unique_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    bulk.flow_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = bulk.flow_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = bulk.flow_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    bulk.flow_unique_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    bulk.flow_unique_attribute_values.flowID = $2::uuid
     AND    $3::timestamptz(3) <@ bulk.flow_unique_attribute_values.validity
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            bulk.flow_standard_attribute_values.value AS value,
            lower(bulk.flow_standard_attribute_values.validity) AS validSince,
            upper(bulk.flow_standard_attribute_values.validity) AS validUntil,
            bulk.flow_standard_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    bulk.flow_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = bulk.flow_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = bulk.flow_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    bulk.flow_standard_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    bulk.flow_standard_attribute_values.flowID = $2::uuid
     AND    $3::timestamptz(3) <@ bulk.flow_standard_attribute_values.validity;`
)

func init() {
	m[FlowList] = `FlowList`
	m[FlowTxShowProperties] = `FlowTxShowProperties`
	m[FlowTxShow] = `FlowTxShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
