/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	ContainerReadStatements = ``

	ContainerList = `
SELECT            meta.dictionary.name AS dictionaryName,
                  asset.container_unique_attribute_values.value AS runtimeName,
                  inventory.user.uid AS createdBy,
                  asset.container_unique_attribute_values.createdAt
FROM              meta.dictionary
JOIN              meta.unique_attribute
  ON              meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
JOIN              asset.container
  ON              meta.dictionary.dictionaryID = asset.container.dictionaryID
JOIN              asset.container_unique_attribute_values
    ON            meta.dictionary.dictionaryID = asset.container_unique_attribute_values.dictionaryID
    AND           asset.container.containerID = asset.container_unique_attribute_values.containerID
    AND           meta.unique_attribute.attributeID = asset.container_unique_attribute_values.attributeID
JOIN              inventory.user
  ON              asset.container_unique_attribute_values.createdBy = inventory.user.userID
WHERE             (meta.dictionary.name = $1::text OR $1::text IS NULL)
  AND             meta.unique_attribute.attribute = 'name'::text
  AND             now()::timestamptz(3) <@ asset.container_unique_attribute_values.validity;`

	ContainerTxShow = `
SELECT            asset.container.containerID,
                  asset.container.dictionaryID,
                  asset.container.createdAt,
                  creator.uid AS createdBy,
                  lower(asset.container_unique_attribute_values.validity) AS validSince,
                  upper(asset.container_unique_attribute_values.validity) AS validUntil,
                  asset.container_unique_attribute_values.createdAt AS namedAt,
                  namegiver.uid AS namedBy
FROM              meta.dictionary
    JOIN          asset.container
        ON        meta.dictionary.dictionaryID = asset.container.dictionaryID
    JOIN          inventory.user AS creator
        ON        asset.container.createdBy = creator.userID
    JOIN          meta.unique_attribute
        ON        meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN          asset.container_unique_attribute_values
        ON        meta.dictionary.dictionaryID = asset.container_unique_attribute_values.dictionaryID
        AND       asset.container.containerID = asset.container_unique_attribute_values.containerID
        AND       meta.unique_attribute.attributeID = asset.container_unique_attribute_values.attributeID
    JOIN          inventory.user AS namegiver
        ON        asset.container_unique_attribute_values.createdBy = namegiver.userID
WHERE             meta.dictionary.name = $1::text
     AND          meta.unique_attribute.attribute = 'name'::text
     AND          asset.container_unique_attribute_values.value = $2::text
     AND          $3::timestamptz(3) <@ asset.container_unique_attribute_values.validity;`

	ContainerTxShowProperties = `
SELECT      meta.unique_attribute.attribute AS attribute,
            asset.container_unique_attribute_values.value AS value,
            lower(asset.container_unique_attribute_values.validity) AS validSince,
            upper(asset.container_unique_attribute_values.validity) AS validUntil,
            asset.container_unique_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    asset.container_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = asset.container_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = asset.container_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.container_unique_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.container_unique_attribute_values.containerID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.container_unique_attribute_values.validity
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            asset.container_standard_attribute_values.value AS value,
            lower(asset.container_standard_attribute_values.validity) AS validSince,
            upper(asset.container_standard_attribute_values.validity) AS validUntil,
            asset.container_standard_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    asset.container_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = asset.container_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = asset.container_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    asset.container_standard_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    asset.container_standard_attribute_values.containerID = $2::uuid
     AND    $3::timestamptz(3) <@ asset.container_standard_attribute_values.validity;`
)

func init() {
	m[ContainerList] = `ContainerList`
	m[ContainerTxShow] = `ContainerTxShow`
	m[ContainerTxShowProperties] = `ContainerTxShowProperties`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
