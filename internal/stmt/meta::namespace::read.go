/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	NamespaceReadStatements = ``

	NamespaceList = `
SELECT      meta.dictionary.dictionaryID,
            meta.unique_attribute.attribute AS key,
            meta.dictionary_unique_attribute_values.value AS value,
            meta.dictionary.createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    meta.dictionary_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = meta.dictionary_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = meta.dictionary_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    meta.dictionary.createdBy =  inventory.user.userID
WHERE       meta.unique_attribute.attribute = 'dict_name'::text
     AND    now()::timestamptz(3) <@ meta.dictionary_unique_attribute_values.validity
UNION
SELECT      meta.dictionary.dictionaryID,
            meta.standard_attribute.attribute AS key,
            meta.dictionary_standard_attribute_values.value AS value,
            meta.dictionary.createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    meta.dictionary_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = meta.dictionary_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    meta.dictionary.createdBy =  inventory.user.userID
WHERE       meta.standard_attribute.attribute = 'dict_type'::text
     AND    now()::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity;`

	NamespaceTxShow = `
SELECT      meta.dictionary.dictionaryID,
            meta.dictionary.name,
            meta.dictionary.createdAt,
            inventory.user.uid
FROM        meta.dictionary
    JOIN    inventory.user
      ON    meta.dictionary.createdBy = inventory.user.userID
WHERE       meta.dictionary.name = $1::text;`

	NamespaceTxSelectProperties = `
SELECT      meta.unique_attribute.attribute AS attribute,
            meta.dictionary_unique_attribute_values.value AS value,
            lower(meta.dictionary_unique_attribute_values.validity) AS validSince,
            upper(meta.dictionary_unique_attribute_values.validity) AS validUntil,
            meta.dictionary_unique_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    meta.dictionary_unique_attribute_values
      ON    meta.unique_attribute.dictionaryID = meta.dictionary_unique_attribute_values.dictionaryID
     AND    meta.unique_attribute.attributeID  = meta.dictionary_unique_attribute_values.attributeID
    JOIN    inventory.user
      ON    meta.dictionary_unique_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    $2::timestamptz(3) <@ meta.dictionary_unique_attribute_values.validity
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            meta.dictionary_standard_attribute_values.value AS value,
            lower(meta.dictionary_standard_attribute_values.validity) AS validSince,
            upper(meta.dictionary_standard_attribute_values.validity) AS validUntil,
            meta.dictionary_standard_attribute_values.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    meta.dictionary_standard_attribute_values
      ON    meta.standard_attribute.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
     AND    meta.standard_attribute.attributeID  = meta.dictionary_standard_attribute_values.attributeID
    JOIN    inventory.user
      ON    meta.dictionary_standard_attribute_values.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
     AND    $2::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity;`

	NamespaceTxSelectAttributes = `
SELECT      meta.unique_attribute.attribute AS attribute,
            'unique'::text AS type,
            meta.unique_attribute.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.unique_attribute
      ON    meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN    inventory.user
      ON    meta.unique_attribute.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid
UNION
SELECT      meta.standard_attribute.attribute AS attribute,
            'standard'::text AS type,
            meta.standard_attribute.createdAt AS createdAt,
            inventory.user.uid AS uid
FROM        meta.dictionary
    JOIN    meta.standard_attribute
      ON    meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
    JOIN    inventory.user
      ON    meta.standard_attribute.createdBy =  inventory.user.userID
WHERE       meta.dictionary.dictionaryID = $1::uuid;`

	NamespaceAttributeSelect = `
SELECT            meta.standard_attribute.attributeID,
                  meta.standard_attribute.dictionaryID,
                  'standard'::text AS attributeType
FROM              meta.dictionary
    JOIN          meta.attribute
      ON          meta.dictionary.dictionaryID = meta.attribute.dictionaryID
    JOIN          meta.standard_attribute
      ON          meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
     AND          meta.attribute.attribute = meta.standard_attribute.attribute
WHERE             meta.dictionary.name = $1::text
  AND             meta.attribute.attribute = $2::text
UNION
SELECT            meta.unique_attribute.attributeID,
                  meta.unique_attribute.dictionaryID,
                  'unique'::text AS attributeType
FROM              meta.dictionary
    JOIN          meta.attribute
      ON          meta.dictionary.dictionaryID = meta.attribute.dictionaryID
    JOIN          meta.unique_attribute
      ON          meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
     AND          meta.attribute.attribute = meta.unique_attribute.attribute
WHERE             meta.dictionary.name = $1::text
  AND             meta.attribute.attribute = $2::text;`
)

func init() {
	m[NamespaceAttributeSelect] = `NamespaceAttributeSelect`
	m[NamespaceList] = `NamespaceList`
	m[NamespaceTxSelectAttributes] = `NamespaceTxSelectAttributes`
	m[NamespaceTxSelectProperties] = `NamespaceTxSelectProperties`
	m[NamespaceTxShow] = `NamespaceTxShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
