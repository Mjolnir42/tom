/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	NamespaceWriteStatements = ``

	NamespaceAdd = `
WITH ins_dct AS ( INSERT INTO meta.dictionary ( name, createdBy )
                  VALUES      ( $1::text, ( SELECT inventory.user.userID
                                FROM inventory.user
                                JOIN inventory.identity_library
                                  ON inventory.identity_library.identityLibraryID
                                   = inventory.user.identityLibraryID
                                WHERE inventory.user.uid = $3::text
                                  AND inventory.identity_library.name = $2::text))
                  ON CONFLICT ( name ) DO UPDATE SET name=EXCLUDED.name
                  RETURNING   dictionaryID AS dictID, name AS dictName, createdBy AS userID),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute, createdBy )
                  SELECT      dictID,
                              'dict_name',
                              userID
                  FROM        ins_dct ),
     ins_nam AS ( INSERT INTO meta.unique_attribute ( dictionaryID, attribute, createdBy )
                  SELECT      dictID,
                              'dict_name',
                              userID
                  FROM        ins_dct
                  ON CONFLICT ON CONSTRAINT __uniq_unique_attr DO UPDATE SET dictionaryID=EXCLUDED.dictionaryID
                  RETURNING   attributeID AS attrID )
INSERT INTO       meta.dictionary_unique_attribute_values ( dictionaryID, attributeID, value, validity, createdBy )
SELECT            dictID,
                  attrID,
                  dictName,
                  '[-infinity,infinity]'::tstzrange,
                  userID
FROM              ins_dct
  CROSS JOIN      ins_nam
ON CONFLICT       ON CONSTRAINT __mdq_temporal DO NOTHING;`

	NamespaceConfigure = `
WITH cte     AS ( SELECT      meta.dictionary.dictionaryID AS dictID,
                              inventory.user.userID AS userID
                  FROM        meta.dictionary
                  CROSS JOIN  inventory.user
                        JOIN  inventory.identity_library
                          ON  inventory.identity_library.identityLibraryID
                           =  inventory.user.identityLibraryID
                  WHERE       meta.dictionary.name = $1::text
                    AND       inventory.user.uid = $5::text
                    AND       inventory.identity_library.name = $4::text),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute, createdBy )
                  SELECT      dictID,
                              $2::text,
                              userID
                  FROM        cte ),
     ins_typ AS ( INSERT INTO meta.standard_attribute ( dictionaryID, attribute, createdBy )
                  SELECT      dictID,
                              $2::text,
                              userID
                  FROM        cte
                  ON CONFLICT ON CONSTRAINT __uniq_attribute DO UPDATE SET dictionaryID=EXCLUDED.dictionaryID
                  RETURNING   attributeID AS attrID )
INSERT INTO       meta.dictionary_standard_attribute_values ( dictionaryID, attributeID, value, validity, createdBy )
SELECT            dictID,
                  attrID,
                  $3::text,
                  '[-infinity,infinity]'::tstzrange,
                  userID
FROM              cte
  CROSS JOIN      ins_typ
ON CONFLICT       ON CONSTRAINT __mda_temporal DO NOTHING;`

	NamespaceRemove = `
SELECT      'Namespace.REMOVE(TODO)';`

	NamespaceAttributeAddStandard = `
WITH cte     AS ( SELECT      meta.dictionary.dictionaryID AS dictID,
                              inventory.user.userID AS userID
                  FROM        meta.dictionary
                  CROSS JOIN  inventory.user
                        JOIN  inventory.identity_library
                          ON  inventory.identity_library.identityLibraryID
                           =  inventory.user.identityLibraryID
                  WHERE       meta.dictionary.name = $1::text
                    AND       inventory.user.uid = $4::text
                    AND       inventory.identity_library.name = $3::text),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute, createdBy )
                  SELECT      dictID,
                              $2::text,
                              userID
                  FROM        cte )
INSERT INTO       meta.standard_attribute ( dictionaryID, attribute, createdBy )
SELECT            cte.dictID,
                  v.value,
                  cte.userID
FROM              cte
  CROSS JOIN      (VALUES($2::text)) AS v (value)
ON CONFLICT       ON CONSTRAINT __uniq_attribute DO NOTHING;`

	NamespaceAttributeAddUnique = `
WITH cte     AS ( SELECT      meta.dictionary.dictionaryID AS dictID,
                              inventory.user.userID AS userID
                  FROM        meta.dictionary
                  CROSS JOIN  inventory.user
                        JOIN  inventory.identity_library
                          ON  inventory.identity_library.identityLibraryID
                           =  inventory.user.identityLibraryID
                  WHERE       meta.dictionary.name = $1::text
                    AND       inventory.user.uid = $4::text
                    AND       inventory.identity_library.name = $3::text),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute, createdBy )
                  SELECT      dictID,
                              $2::text,
                              userID
                  FROM        cte )
INSERT INTO       meta.unique_attribute ( dictionaryID, attribute, createdBy )
SELECT            cte.dictID,
                  v.value,
                  cte.userID
FROM              cte
  CROSS JOIN      (VALUES($2::text)) AS v (value)
ON CONFLICT       ON CONSTRAINT __uniq_unique_attr DO NOTHING;`

	NamespaceAttributeQueryType = `
SELECT            'standard'::text AS attributeType
FROM              meta.dictionary
    JOIN          meta.attribute
      ON          meta.dictionary.dictionaryID = meta.attribute.dictionaryID
    JOIN          meta.standard_attribute
      ON          meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
     AND          meta.attribute.attribute = meta.standard_attribute.attribute
WHERE             meta.dictionary.name = $1::text
  AND             meta.attribute.attribute = $2::text
UNION
SELECT            'unique'::text AS attributeType
FROM              meta.dictionary
    JOIN          meta.attribute
      ON          meta.dictionary.dictionaryID = meta.attribute.dictionaryID
    JOIN          meta.unique_attribute
      ON          meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
     AND          meta.attribute.attribute = meta.unique_attribute.attribute
WHERE             meta.dictionary.name = $1::text
  AND             meta.attribute.attribute = $2::text;`

	NamespaceAttributeDiscover = `
SELECT            meta.attribute.attribute AS attributeName,
                  'standard'::text AS attributeType
FROM              meta.dictionary
    JOIN          meta.attribute
      ON          meta.dictionary.dictionaryID = meta.attribute.dictionaryID
    JOIN          meta.standard_attribute
      ON          meta.dictionary.dictionaryID = meta.standard_attribute.dictionaryID
     AND          meta.attribute.attribute = meta.standard_attribute.attribute
WHERE             meta.dictionary.name = $1::text
  AND             meta.attribute.attribute NOT LIKE 'dict_%'
UNION
SELECT            meta.attribute.attribute AS attributeName,
                  'unique'::text AS attributeType
FROM              meta.dictionary
    JOIN          meta.attribute
      ON          meta.dictionary.dictionaryID = meta.attribute.dictionaryID
    JOIN          meta.unique_attribute
      ON          meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
     AND          meta.attribute.attribute = meta.unique_attribute.attribute
WHERE             meta.dictionary.name = $1::text
  AND             meta.attribute.attribute NOT LIKE 'dict_%';`

	NamespaceTxStdPropertySelect = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
SELECT            value,
                  lower(validity),
                  upper(validity)
FROM              meta.dictionary_standard_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = meta.dictionary_standard_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity
   FOR UPDATE;`

	NamespaceTxStdPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
UPDATE            meta.dictionary_standard_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             meta.dictionary_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             meta.dictionary_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             meta.dictionary_standard_attribute_values.value           = $3::text
  AND             lower(meta.dictionary_standard_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(meta.dictionary_standard_attribute_values.validity) = $5::timestamptz(3)
  AND             $6::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity;`

	NamespaceTxStdPropertyAdd = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID AS dictID,
                              inventory.user.userID AS userID
                  FROM        meta.dictionary
                  CROSS JOIN  inventory.user
                        JOIN  inventory.identity_library
                          ON  inventory.identity_library.identityLibraryID
                           =  inventory.user.identityLibraryID
                  WHERE       meta.dictionary.name = $1::text
                    AND       inventory.user.uid = $7::text
                    AND       inventory.identity_library.name = $6::text),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
INSERT INTO       meta.dictionary_standard_attribute_values ( dictionaryID, attributeID, value, validity, createdBy )
SELECT            cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(4), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	NamespaceTxUniqPropertySelect = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
SELECT            value,
                  lower(validity),
                  upper(validity)
FROM              meta.dictionary_unique_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = meta.dictionary_unique_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = meta.dictionary_unique_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ meta.dictionary_unique_attribute_values.validity
   FOR UPDATE;`

	NamespaceTxUniqPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
UPDATE            meta.dictionary_unique_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             meta.dictionary_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             meta.dictionary_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             meta.dictionary_unique_attribute_values.value           = $3::text
  AND             lower(meta.dictionary_unique_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(meta.dictionary_unique_attribute_values.validity) = $5::timestamptz(3)
  AND             $6::timestamptz(3) <@ meta.dictionary_unique_attribute_values.validity;`

	NamespaceTxUniqPropertyAdd = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID AS dictID,
                              inventory.user.userID AS userID
                  FROM        meta.dictionary
                  CROSS JOIN  inventory.user
                        JOIN  inventory.identity_library
                          ON  inventory.identity_library.identityLibraryID
                           =  inventory.user.identityLibraryID
                  WHERE       meta.dictionary.name = $1::text
                    AND       inventory.user.uid = $7::text
                    AND       inventory.identity_library.name = $6::text),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                  FROM        meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
INSERT INTO       meta.dictionary_unique_attribute_values ( dictionaryID, attributeID, value, validity, createdBy )
SELECT            cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`
)

func init() {
	m[NamespaceAdd] = `NamespaceAdd`
	m[NamespaceAttributeAddStandard] = `NamespaceAttributeAddStandard`
	m[NamespaceAttributeAddUnique] = `NamespaceAttributeAddUnique`
	m[NamespaceAttributeDiscover] = `NamespaceAttributeDiscover`
	m[NamespaceAttributeQueryType] = `NamespaceAttributeQueryType`
	m[NamespaceConfigure] = `NamespaceConfigure`
	m[NamespaceRemove] = `NamespaceRemove`
	m[NamespaceTxStdPropertyAdd] = `NamespaceTxStdPropertyAdd`
	m[NamespaceTxStdPropertyClamp] = `NamespaceTxStdPropertyClamp`
	m[NamespaceTxStdPropertySelect] = `NamespaceTxStdPropertySelect`
	m[NamespaceTxUniqPropertyAdd] = `NamespaceTxUniqPropertyAdd`
	m[NamespaceTxUniqPropertyClamp] = `NamespaceTxUniqPropertyClamp`
	m[NamespaceTxUniqPropertySelect] = `NamespaceTxUniqPropertySelect`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
