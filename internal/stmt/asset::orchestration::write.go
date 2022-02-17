/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	OrchestrationWriteStatements = ``

	OrchestrationAdd = `
WITH sel_dct AS ( SELECT dictionaryID
                  FROM   meta.dictionary
                  WHERE  name = $1::text),
     ins_orc AS ( INSERT INTO asset.orchestration_environment ( dictionaryID, createdBy )
                  VALUES      (( SELECT dictionaryID FROM sel_dct ),
                               ( SELECT inventory.user.userID
                                 FROM inventory.user
                                 JOIN inventory.identity_library
                                 ON inventory.identity_library.identityLibraryID
                                  = inventory.user.identityLibraryID
                                 WHERE inventory.user.uid = $3::text
                                   AND inventory.identity_library.name = $2::text))
                  RETURNING orchID, createdBy AS userID ),
     sel_att AS ( SELECT attributeID
                  FROM   meta.unique_attribute
                  JOIN   meta.dictionary
                    ON   meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE  meta.dictionary.name = $1::text
                    AND  meta.unique_attribute.attribute = 'name'::text )
INSERT INTO       asset.orchestration_environment_unique_attribute_values (
                         orchID,
                         attributeID,
                         dictionaryID,
                         value,
                         validity,
                         createdBy
                  )
SELECT            ins_orc.orchID,
                  sel_att.attributeID,
                  sel_dct.dictionaryID,
                  $4::text,
                  tstzrange( $5::timestamptz(3), $6::timestamptz(3), '[]'),
                  ins_orc.userID
FROM              ins_orc
  CROSS JOIN      sel_att
  CROSS JOIN      sel_dct
RETURNING         orchID;`

	OrchestrationTxStdPropertyAdd = `
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
INSERT INTO       asset.orchestration_environment_standard_attribute_values ( orchID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	OrchestrationTxStdPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
UPDATE            asset.orchestration_environment_standard_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.orchestration_environment_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.orchestration_environment_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.orchestration_environment_standard_attribute_values.value           = $3::text
  AND             lower(asset.orchestration_environment_standard_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.orchestration_environment_standard_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.orchestration_environment_standard_attribute_values.orchID          = $8::uuid
  AND             $6::timestamptz(3) <@ asset.orchestration_environment_standard_attribute_values.validity;`

	OrchestrationTxStdPropertyClean = `
-- this statement is used to delete all records, for which the starting validity is after the timestamp specified
-- in $4. this can be used to clean all records that only become valid in the future
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
DELETE FROM       asset.orchestration_environment_standard_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.orchestration_environment_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.orchestration_environment_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.orchestration_environment_standard_attribute_values.orchID          = $3::uuid
  AND             lower(asset.orchestration_environment_standard_attribute_values.validity) > $4::timestamptz(3);`

	OrchestrationTxStdPropertySelect = `
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
FROM              asset.orchestration_environment_standard_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.orchestration_environment_standard_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.orchestration_environment_standard_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.orchestration_environment_standard_attribute_values.validity
     AND          asset.orchestration_environment_standard_attribute_values.orchID = $4::uuid
   FOR UPDATE;`

	OrchestrationTxUniqPropertyAdd = `
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
INSERT INTO       asset.orchestration_environment_unique_attribute_values ( orchID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	OrchestrationTxUniqPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
UPDATE            asset.orchestration_environment_unique_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.orchestration_environment_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.orchestration_environment_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.orchestration_environment_unique_attribute_values.value           = $3::text
  AND             lower(asset.orchestration_environment_unique_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.orchestration_environment_unique_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.orchestration_environment_unique_attribute_values.orchID        = $8::uuid
  AND             $6::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity;`

	OrchestrationTxUniqPropertySelect = `
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
FROM              asset.orchestration_environment_unique_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.orchestration_environment_unique_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.orchestration_environment_unique_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity
     AND          asset.orchestration_environment_unique_attribute_values.orchID = $4::uuid
   FOR UPDATE;`

	OrchestrationTxUniqPropertyClean = `
-- this statement is used to delete all records, for which the starting validity is after the timestamp specified
-- in $4. this can be used to clean all records that only become valid in the future
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
DELETE FROM       asset.orchestration_environment_unique_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.orchestration_environment_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.orchestration_environment_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.orchestration_environment_unique_attribute_values.orchID          = $3::uuid
  AND             lower(asset.orchestration_environment_unique_attribute_values.validity) > $4::timestamptz(3);`

	OrchestrationTxStackAdd = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $6::text
                    AND  inventory.identity_library.name = $7::text)
INSERT INTO       asset.orchestration_environment_mapping (
                         orchID,
                         parentRuntimeID,
                         validity,
                         createdBy,
                         createdAt
                  )
SELECT            $1::uuid,
                  $2::uuid,
                  tstzrange($3::timestamptz(3), $4::timestamptz(3), '[]'),
                  sel_uid.userID,
                  $5::timestamptz(3)
FROM              sel_uid;`

	OrchestrationTxStackClamp = `
UPDATE            asset.orchestration_environment_mapping
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.orchestration_environment_mapping.orchID = $2::uuid
  AND             asset.orchestration_environment_mapping.parentRuntimeID = $3::uuid
  AND             $1::timestamptz(3) <@ asset.orchestration_environment_mapping.validity
	AND             $4::timestamptz(3) <@ asset.orchestration_environment_mapping.validity;`

	OrchestrationTxStackClampAll = `
UPDATE            asset.orchestration_environment_mapping
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.orchestration_environment_mapping.orchID = $2::uuid
	AND             $1::timestamptz(3) <@ asset.orchestration_environment_mapping.validity;`

	OrchestrationTxLink = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $5::text
                    AND  inventory.identity_library.name = $6::text)
INSERT INTO       asset.orchestration_environment_linking (
                         orchID_A,
                         dictionaryID_A,
                         orchID_B,
                         dictionaryID_B,
                         createdBy,
                         createdAt
                  )
SELECT            $1::uuid,
                  $2::uuid,
                  $3::uuid,
                  $4::uuid,
                  sel_uid.userID,
                  $7::timestamptz(3)
FROM              sel_uid;`

	OrchestrationTxUnstackChild = `
UPDATE            asset.runtime_environment_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.runtime_environment_parent.parentOrchestrationID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.runtime_environment_parent.validity;`

	OrchestrationTxUnstackChildClean = `
DELETE FROM       asset.runtime_environment_parent
WHERE             asset.runtime_environment_parent.parentOrchestrationID = $2::uuid
  AND             lower(asset.runtime_environment_parent.validity) > $1::timestamptz(3);`

	OrchestrationDelNamespaceStdValues = `
DELETE FROM       asset.orchestration_environment_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	OrchestrationDelNamespaceUniqValues = `
DELETE FROM       asset.orchestration_environment_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	OrchestrationDelNamespace = `
DELETE FROM       asset.orchestration_environment
WHERE             dictionaryID = $1::uuid;`

	OrchestrationDelNamespaceLinking = `
DELETE FROM       asset.orchestration_environment_linking
WHERE             dictionaryID_A = $1::uuid
   OR             dictionaryID_B = $1::uuid;`

	OrchestrationDelNamespaceParent = `
DELETE FROM       asset.orchestration_environment_mapping
USING             asset.orchestration_environment
WHERE             asset.orchestration_environment_mapping.orchID = asset.orchestration_environment.orchID
  AND             asset.orchestration_environment.dictionaryID = $1::uuid;`
)

func init() {
	m[OrchestrationAdd] = `OrchestrationAdd`
	m[OrchestrationDelNamespaceLinking] = `OrchestrationDelNamespaceLinking`
	m[OrchestrationDelNamespaceParent] = `OrchestrationDelNamespaceParent`
	m[OrchestrationDelNamespaceStdValues] = `OrchestrationDelNamespaceStdValues`
	m[OrchestrationDelNamespaceUniqValues] = `OrchestrationDelNamespaceUniqValues`
	m[OrchestrationDelNamespace] = `OrchestrationDelNamespace`
	m[OrchestrationTxLink] = `OrchestrationTxLink`
	m[OrchestrationTxStackAdd] = `OrchestrationTxStackAdd`
	m[OrchestrationTxStackClampAll] = `OrchestrationTxStackClampAll`
	m[OrchestrationTxStackClamp] = `OrchestrationTxStackClamp`
	m[OrchestrationTxStdPropertyAdd] = `OrchestrationTxStdPropertyAdd`
	m[OrchestrationTxStdPropertyClamp] = `OrchestrationTxStdPropertyClamp`
	m[OrchestrationTxStdPropertyClean] = `OrchestrationTxStdPropertyClean`
	m[OrchestrationTxStdPropertySelect] = `OrchestrationTxStdPropertySelect`
	m[OrchestrationTxUniqPropertyAdd] = `OrchestrationTxUniqPropertyAdd`
	m[OrchestrationTxUniqPropertyClamp] = `OrchestrationTxUniqPropertyClamp`
	m[OrchestrationTxUniqPropertyClean] = `OrchestrationTxUniqPropertyClean`
	m[OrchestrationTxUniqPropertySelect] = `OrchestrationTxUniqPropertySelect`
	m[OrchestrationTxUnstackChildClean] = `OrchestrationTxUnstackChildClean`
	m[OrchestrationTxUnstackChild] = `OrchestrationTxUnstackChild`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
