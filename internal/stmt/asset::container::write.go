/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	ContainerWriteStatements = ``

	ContainerAdd = `
WITH sel_dct AS ( SELECT dictionaryID
                  FROM   meta.dictionary
                  WHERE  name = $1::text),
     ins_rte AS ( INSERT INTO asset.container ( dictionaryID, createdBy )
                  VALUES      (( SELECT dictionaryID FROM sel_dct ),
                               ( SELECT inventory.user.userID
                                 FROM inventory.user
                                 JOIN inventory.identity_library
                                 ON inventory.identity_library.identityLibraryID
                                  = inventory.user.identityLibraryID
                                 WHERE inventory.user.uid = $3::text
                                   AND inventory.identity_library.name = $2::text))
                  RETURNING containerID, createdBy AS userID ),
     sel_att AS ( SELECT attributeID
                  FROM   meta.unique_attribute
                  JOIN   meta.dictionary
                    ON   meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE  meta.dictionary.name = $1::text
                    AND  meta.unique_attribute.attribute = 'name'::text )
INSERT INTO       asset.container_unique_attribute_values (
                         containerID,
                         attributeID,
                         dictionaryID,
                         value,
                         validity,
                         createdBy
                  )
SELECT            ins_rte.containerID,
                  sel_att.attributeID,
                  sel_dct.dictionaryID,
                  $4::text,
                  tstzrange( $5::timestamptz(3), $6::timestamptz(3), '[]'),
                  ins_rte.userID
FROM              ins_rte
  CROSS JOIN      sel_att
  CROSS JOIN      sel_dct
RETURNING         containerID;`

	ContainerRemove = `
SELECT      'Container.REMOVE';`

	ContainerTxStdPropertyAdd = `
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
INSERT INTO       asset.container_standard_attribute_values ( containerID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	ContainerTxStdPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
UPDATE            asset.container_standard_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.container_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.container_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.container_standard_attribute_values.value           = $3::text
  AND             lower(asset.container_standard_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.container_standard_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.container_standard_attribute_values.containerID     = $8::uuid
  AND             $6::timestamptz(3) <@ asset.container_standard_attribute_values.validity;`

	ContainerTxStdPropertySelect = `
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
FROM              asset.container_standard_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.container_standard_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.container_standard_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.container_standard_attribute_values.validity
     AND          asset.container_standard_attribute_values.containerID = $4::uuid
   FOR UPDATE;`

	ContainerTxUniqPropertyAdd = `
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
INSERT INTO       asset.container_unique_attribute_values ( containerID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	ContainerTxUniqPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
UPDATE            asset.container_unique_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.container_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.container_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.container_unique_attribute_values.value           = $3::text
  AND             lower(asset.container_unique_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.container_unique_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.container_unique_attribute_values.containerID     = $8::uuid
  AND             $6::timestamptz(3) <@ asset.container_unique_attribute_values.validity;`

	ContainerTxUniqPropertySelect = `
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
FROM              asset.container_unique_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.container_unique_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.container_unique_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.container_unique_attribute_values.validity
     AND          asset.container_unique_attribute_values.containerID = $4::uuid
   FOR UPDATE;`

	ContainerDelNamespaceStdValues = `
DELETE FROM       asset.container_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	ContainerDelNamespaceUniqValues = `
DELETE FROM       asset.container_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	ContainerDelNamespace = `
DELETE FROM       asset.container
WHERE             dictionaryID = $1::uuid;`

	ContainerDelNamespaceLinking = `
DELETE FROM       asset.container_linking
WHERE             dictionaryID_A = $1::uuid
   OR             dictionaryID_B = $1::uuid;`

	ContainerDelNamespaceParent = `
DELETE FROM       asset.container_parent
USING             asset.container
WHERE             asset.container_parent.containerID = asset.container.containerID
  AND             asset.container.dictionaryID = $1::uuid;`

	ContainerLink = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $5::text
                    AND  inventory.identity_library.name = $6::text)
INSERT INTO       asset.container_linking (
                         containerID_A,
                         dictionaryID_A,
                         containerID_B,
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

	ContainerTxStackAdd = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $6::text
                    AND  inventory.identity_library.name = $7::text)
INSERT INTO       asset.container_parent (
                         containerID,
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

	ContainerTxStackClamp = `
UPDATE            asset.container_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.container_parent.containerID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.container_parent.validity;`

	ContainerTxStdPropertyClean = `
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
DELETE FROM       asset.container_standard_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.container_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.container_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.container_standard_attribute_values.containerID     = $3::uuid
  AND             lower(asset.container_standard_attribute_values.validity) > $4::timestamptz(3);`

	ContainerTxUniqPropertyClean = `
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
DELETE FROM       asset.container_unique_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.container_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.container_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.container_unique_attribute_values.containerID     = $3::uuid
  AND             lower(asset.container_unique_attribute_values.validity) > $4::timestamptz(3);`
)

func init() {
	m[ContainerAdd] = `ContainerAdd`
	m[ContainerDelNamespaceLinking] = `ContainerDelNamespaceLinking`
	m[ContainerDelNamespaceParent] = `ContainerDelNamespaceParent`
	m[ContainerDelNamespaceStdValues] = `ContainerDelNamespaceStdValues`
	m[ContainerDelNamespaceUniqValues] = `ContainerDelNamespaceUniqValues`
	m[ContainerDelNamespace] = `ContainerDelNamespace`
	m[ContainerLink] = `ContainerLink`
	m[ContainerRemove] = `ContainerRemove`
	m[ContainerTxStackAdd] = `ContainerTxStackAdd`
	m[ContainerTxStackClamp] = `ContainerTxStackClamp`
	m[ContainerTxStdPropertyAdd] = `ContainerTxStdPropertyAdd`
	m[ContainerTxStdPropertyClamp] = `ContainerTxStdPropertyClamp`
	m[ContainerTxStdPropertyClean] = `ContainerTxStdPropertyClean`
	m[ContainerTxStdPropertySelect] = `ContainerTxStdPropertySelect`
	m[ContainerTxUniqPropertyAdd] = `ContainerTxUniqPropertyAdd`
	m[ContainerTxUniqPropertyClamp] = `ContainerTxUniqPropertyClamp`
	m[ContainerTxUniqPropertyClean] = `ContainerTxUniqPropertyClean`
	m[ContainerTxUniqPropertySelect] = `ContainerTxUniqPropertySelect`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
