/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	RuntimeWriteStatements = ``

	RuntimeAdd = `
WITH sel_dct AS ( SELECT dictionaryID
                  FROM   meta.dictionary
                  WHERE  name = $1::text),
     ins_rte AS ( INSERT INTO asset.runtime_environment ( dictionaryID, createdBy )
                  VALUES      (( SELECT dictionaryID FROM sel_dct ),
                               ( SELECT inventory.user.userID
                                 FROM inventory.user
                                 JOIN inventory.identity_library
                                 ON inventory.identity_library.identityLibraryID
                                  = inventory.user.identityLibraryID
                                 WHERE inventory.user.uid = $3::text
                                   AND inventory.identity_library.name = $2::text))
                  RETURNING rteID, createdBy AS userID ),
     sel_att AS ( SELECT attributeID
                  FROM   meta.unique_attribute
                  JOIN   meta.dictionary
                    ON   meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE  meta.dictionary.name = $1::text
                    AND  meta.unique_attribute.attribute = 'name'::text )
INSERT INTO       asset.runtime_environment_unique_attribute_values (
                         rteID,
                         attributeID,
                         dictionaryID,
                         value,
                         validity,
                         createdBy
                  )
SELECT            ins_rte.rteID,
                  sel_att.attributeID,
                  sel_dct.dictionaryID,
                  $4::text,
                  tstzrange( $5::timestamptz(3), $6::timestamptz(3), '[]'),
                  ins_rte.userID
FROM              ins_rte
  CROSS JOIN      sel_att
  CROSS JOIN      sel_dct
RETURNING         rteID;`

	RuntimeRemove = `
SELECT      'Runtime.REMOVE';`

	RuntimeStdAttrRemove = `
DELETE FROM       asset.runtime_environment_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	RuntimeUniqAttrRemove = `
DELETE FROM       asset.runtime_environment_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	RuntimeTxStdPropertyAdd = `
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
INSERT INTO       asset.runtime_environment_standard_attribute_values ( rteID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	RuntimeTxStdPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
UPDATE            asset.runtime_environment_standard_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.runtime_environment_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.runtime_environment_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.runtime_environment_standard_attribute_values.value           = $3::text
  AND             lower(asset.runtime_environment_standard_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.runtime_environment_standard_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.runtime_environment_standard_attribute_values.rteID           = $8::uuid
  AND             $6::timestamptz(3) <@ asset.runtime_environment_standard_attribute_values.validity;`

	RuntimeTxStdPropertySelect = `
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
FROM              asset.runtime_environment_standard_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.runtime_environment_standard_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.runtime_environment_standard_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.runtime_environment_standard_attribute_values.validity
     AND          asset.runtime_environment_standard_attribute_values.rteID = $4::uuid
   FOR UPDATE;`

	RuntimeTxUniqPropertyAdd = `
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
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
INSERT INTO       asset.runtime_environment_unique_attribute_values ( rteID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	RuntimeTxUniqPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
UPDATE            asset.runtime_environment_unique_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.runtime_environment_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.runtime_environment_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.runtime_environment_unique_attribute_values.value           = $3::text
  AND             lower(asset.runtime_environment_unique_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.runtime_environment_unique_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.runtime_environment_unique_attribute_values.rteID           = $8::uuid
  AND             $6::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	RuntimeTxUniqPropertySelect = `
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
FROM              asset.runtime_environment_unique_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.runtime_environment_unique_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity
     AND          asset.runtime_environment_unique_attribute_values.rteID = $4::uuid
   FOR UPDATE;`
)

func init() {
	m[RuntimeAdd] = `RuntimeAdd`
	m[RuntimeRemove] = `RuntimeRemove`
	m[RuntimeStdAttrRemove] = `RuntimeStdAttrRemove`
	m[RuntimeTxStdPropertyAdd] = `RuntimeTxStdPropertyAdd`
	m[RuntimeTxStdPropertyClamp] = `RuntimeTxStdPropertyClamp`
	m[RuntimeTxStdPropertySelect] = `RuntimeTxStdPropertySelect`
	m[RuntimeTxUniqPropertyAdd] = `RuntimeTxUniqPropertyAdd`
	m[RuntimeTxUniqPropertyClamp] = `RuntimeTxUniqPropertyClamp`
	m[RuntimeTxUniqPropertySelect] = `RuntimeTxUniqPropertySelect`
	m[RuntimeUniqAttrRemove] = `RuntimeUniqAttrRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
