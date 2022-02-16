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

	RuntimeDelNamespaceStdValues = `
DELETE FROM       asset.runtime_environment_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	RuntimeDelNamespaceUniqValues = `
DELETE FROM       asset.runtime_environment_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	RuntimeDelNamespace = `
DELETE FROM       asset.runtime_environment
WHERE             dictionaryID = $1::uuid;`

	RuntimeDelNamespaceLinking = `
DELETE FROM       asset.runtime_environment_linking
WHERE             dictionaryID_A = $1::uuid
   OR             dictionaryID_B = $1::uuid;`

	RuntimeDelNamespaceParent = `
DELETE FROM       asset.runtime_environment_parent
USING             asset.runtime_environment
WHERE             asset.runtime_environment_parent.rteID = asset.runtime_environment.rteID
  AND             asset.runtime_environment.dictionaryID = $1::uuid;`

	RuntimeLink = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $5::text
                    AND  inventory.identity_library.name = $6::text)
INSERT INTO       asset.runtime_environment_linking (
                         rteID_A,
                         dictionaryID_A,
                         rteID_B,
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

	RuntimeTxStackAdd = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $8::text
                    AND  inventory.identity_library.name = $9::text)
INSERT INTO       asset.runtime_environment_parent (
                         rteID,
                         parentServerID,
                         parentRuntimeID,
                         parentOrchestrationID,
                         validity,
                         createdBy,
                         createdAt
                  )
SELECT            $1::uuid,
                  $2::uuid,
                  $3::uuid,
                  $4::uuid,
                  tstzrange($5::timestamptz(3), $6::timestamptz(3), '[]'),
                  sel_uid.userID,
                  $7::timestamptz(3)
FROM              sel_uid;`

	RuntimeTxStackClamp = `
UPDATE            asset.runtime_environment_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.runtime_environment_parent.rteID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.runtime_environment_parent.validity;`

	RuntimeTxStdPropertyClean = `
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
DELETE FROM       asset.runtime_environment_standard_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.runtime_environment_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.runtime_environment_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.runtime_environment_standard_attribute_values.rteID           = $3::uuid
  AND             lower(asset.runtime_environment_standard_attribute_values.validity) > $4::timestamptz(3);`

	RuntimeTxUniqPropertyClean = `
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
DELETE FROM       asset.runtime_environment_unique_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.runtime_environment_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.runtime_environment_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.runtime_environment_unique_attribute_values.rteID           = $3::uuid
  AND             lower(asset.runtime_environment_unique_attribute_values.validity) > $4::timestamptz(3);`

	RuntimeTxUnstackChildRte = `
UPDATE            asset.runtime_environment_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.runtime_environment_parent.parentRuntimeID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.runtime_environment_parent.validity;`

	RuntimeTxUnstackChildRteClean = `
DELETE FROM       asset.runtime_environment_parent
WHERE             asset.runtime_environment_parent.parentRuntimeID = $2::uuid
  AND             lower(asset.runtime_environment_parent.validity) > $1::timestamptz(3);`

	RuntimeTxUnstackChildSrv = `
UPDATE            asset.server_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.server_parent.parentRuntimeID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.server_parent.validity;`

	RuntimeTxUnstackChildSrvClean = `
DELETE FROM       asset.server_parent
WHERE             asset.server_parent.parentRuntimeID = $2::uuid
  AND             lower(asset.server_parent.validity) > $1::timestamptz(3);`

	RuntimeTxUnstackChildCnr = `
UPDATE            asset.container_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.container_parent.parentRuntimeID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.container_parent.validity;`

	RuntimeTxUnstackChildCnrClean = `
DELETE FROM       asset.container_parent
WHERE             asset.container_parent.parentRuntimeID = $2::uuid
  AND             lower(asset.container_parent.validity) > $1::timestamptz(3);`

	RuntimeTxUnstackChildSok = `
UPDATE            asset.socket_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.socket_parent.parentRuntimeID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.socket_parent.validity;`

	RuntimeTxUnstackChildSokClean = `
DELETE FROM       asset.socket_parent
WHERE             asset.socket_parent.parentRuntimeID = $2::uuid
  AND             lower(asset.socket_parent.validity) > $1::timestamptz(3);`

	RuntimeTxUnstackChildOre = `
UPDATE            asset.orchestration_environment_mapping
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.orchestration_environment_mapping.parentRuntimeID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.orchestration_environment_mapping.validity;`

	RuntimeTxUnstackChildOreClean = `
DELETE FROM       asset.orchestration_environment_mapping
WHERE             asset.orchestration_environment_mapping.parentRuntimeID = $2::uuid
  AND             lower(asset.orchestration_environment_mapping.validity) > $1::timestamptz(3);`
)

func init() {
	m[RuntimeAdd] = `RuntimeAdd`
	m[RuntimeDelNamespaceLinking] = `RuntimeDelNamespaceLinking`
	m[RuntimeDelNamespaceParent] = `RuntimeDelNamespaceParent`
	m[RuntimeDelNamespaceStdValues] = `RuntimeDelNamespaceStdValues`
	m[RuntimeDelNamespaceUniqValues] = `RuntimeDelNamespaceUniqValues`
	m[RuntimeDelNamespace] = `RuntimeDelNamespace`
	m[RuntimeLink] = `RuntimeLink`
	m[RuntimeTxStackAdd] = `RuntimeTxStackAdd`
	m[RuntimeTxStackClamp] = `RuntimeTxStackClamp`
	m[RuntimeTxStdPropertyAdd] = `RuntimeTxStdPropertyAdd`
	m[RuntimeTxStdPropertyClamp] = `RuntimeTxStdPropertyClamp`
	m[RuntimeTxStdPropertyClean] = `RuntimeTxStdPropertyClean`
	m[RuntimeTxStdPropertySelect] = `RuntimeTxStdPropertySelect`
	m[RuntimeTxUniqPropertyAdd] = `RuntimeTxUniqPropertyAdd`
	m[RuntimeTxUniqPropertyClamp] = `RuntimeTxUniqPropertyClamp`
	m[RuntimeTxUniqPropertyClean] = `RuntimeTxUniqPropertyClean`
	m[RuntimeTxUniqPropertySelect] = `RuntimeTxUniqPropertySelect`
	m[RuntimeTxUnstackChildCnrClean] = `RuntimeTxUnstackChildCnrClean`
	m[RuntimeTxUnstackChildCnr] = `RuntimeTxUnstackChildCnr`
	m[RuntimeTxUnstackChildOreClean] = `RuntimeTxUnstackChildOreClean`
	m[RuntimeTxUnstackChildOre] = `RuntimeTxUnstackChildOre`
	m[RuntimeTxUnstackChildRteClean] = `RuntimeTxUnstackChildRteClean`
	m[RuntimeTxUnstackChildRte] = `RuntimeTxUnstackChildRte`
	m[RuntimeTxUnstackChildSokClean] = `RuntimeTxUnstackChildSokClean`
	m[RuntimeTxUnstackChildSok] = `RuntimeTxUnstackChildSok`
	m[RuntimeTxUnstackChildSrvClean] = `RuntimeTxUnstackChildSrvClean`
	m[RuntimeTxUnstackChildSrv] = `RuntimeTxUnstackChildSrv`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
