/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	ServerWriteStatements = ``

	ServerAdd = `
WITH sel_dct AS ( SELECT dictionaryID
                  FROM   meta.dictionary
                  WHERE  name = $1::text),
     ins_srv AS ( INSERT INTO asset.server ( dictionaryID, createdBy )
                  VALUES      (( SELECT dictionaryID FROM sel_dct ),
                               ( SELECT inventory.user.userID
                                 FROM inventory.user
                                 JOIN inventory.identity_library
                                 ON inventory.identity_library.identityLibraryID
                                  = inventory.user.identityLibraryID
                                 WHERE inventory.user.uid = $3::text
                                   AND inventory.identity_library.name = $2::text))
                  RETURNING serverID, createdBy AS userID ),
     sel_att AS ( SELECT attributeID
                  FROM   meta.unique_attribute
                  JOIN   meta.dictionary
                    ON   meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE  meta.dictionary.name = $1::text
                    AND  meta.unique_attribute.attribute = 'name'::text )
INSERT INTO       asset.server_unique_attribute_values (
                         serverID,
                         attributeID,
                         dictionaryID,
                         value,
                         validity,
                         createdBy
                  )
SELECT            ins_srv.serverID,
                  sel_att.attributeID,
                  sel_dct.dictionaryID,
                  $4::text,
                  tstzrange( $5::timestamptz(3), $6::timestamptz(3), '[]'),
                  ins_srv.userID
FROM              ins_srv
  CROSS JOIN      sel_att
	CROSS JOIN      sel_dct
RETURNING         serverID;`

	ServerTxStdPropertyAdd = `
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
INSERT INTO       asset.server_standard_attribute_values ( serverID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	ServerTxStdPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM       meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE       meta.standard_attribute.attribute = $2::text )
UPDATE            asset.server_standard_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.server_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.server_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.server_standard_attribute_values.value           = $3::text
  AND             lower(asset.server_standard_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.server_standard_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.server_standard_attribute_values.serverID        = $8::uuid
  AND             $6::timestamptz(3) <@ asset.server_standard_attribute_values.validity;`

	ServerTxStdPropertySelect = `
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
FROM              asset.server_standard_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.server_standard_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.server_standard_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.server_standard_attribute_values.validity
     AND          asset.server_standard_attribute_values.serverID = $4::uuid
   FOR UPDATE;`

	ServerTxStdPropertyClean = `
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
DELETE FROM       asset.server_standard_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.server_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.server_standard_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.server_standard_attribute_values.serverID        = $3::uuid
  AND             lower(asset.server_standard_attribute_values.validity) > $4::timestamptz(3);`

	ServerTxUniqPropertyAdd = `
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
INSERT INTO       asset.server_unique_attribute_values ( serverID, dictionaryID, attributeID, value, validity, createdBy )
SELECT            $8::uuid,
                  cte_dct.dictID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), $5::timestamptz(3), '[]'),
                  cte_dct.userID
FROM              cte_dct
    CROSS JOIN    cte_att;`

	ServerTxUniqPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM       meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE       meta.unique_attribute.attribute = $2::text )
UPDATE            asset.server_unique_attribute_values
   SET            validity = tstzrange(lower(validity), $7::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE             asset.server_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.server_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.server_unique_attribute_values.value           = $3::text
  AND             lower(asset.server_unique_attribute_values.validity) = $4::timestamptz(3)
  AND             upper(asset.server_unique_attribute_values.validity) = $5::timestamptz(3)
  AND             asset.server_unique_attribute_values.serverID        = $8::uuid
  AND             $6::timestamptz(3) <@ asset.server_unique_attribute_values.validity;`

	ServerTxUniqPropertySelect = `
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
FROM              asset.server_unique_attribute_values
    JOIN          cte_dct
      ON          cte_dct.dictionaryID = asset.server_unique_attribute_values.dictionaryID
    JOIN          cte_att
      ON          cte_att.attributeID = asset.server_unique_attribute_values.attributeID
   WHERE          $3::timestamptz(3) <@ asset.server_unique_attribute_values.validity
     AND          asset.server_unique_attribute_values.serverID = $4::uuid
   FOR UPDATE;`

	ServerTxUniqPropertyClean = `
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
DELETE FROM       asset.server_unique_attribute_values
USING             cte_dct
    CROSS JOIN    cte_att
WHERE             asset.server_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND             asset.server_unique_attribute_values.attributeID     = cte_att.attributeID
  AND             asset.server_unique_attribute_values.serverID        = $3::uuid
  AND             lower(asset.server_unique_attribute_values.validity) > $4::timestamptz(3);`

	ServerLink = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $5::text
                    AND  inventory.identity_library.name = $6::text)
INSERT INTO       asset.server_linking (
                         serverID_A,
                         dictionaryID_A,
                         serverID_B,
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

	ServerTxStackAdd = `
WITH sel_uid AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                    =    inventory.user.identityLibraryID
                  WHERE  inventory.user.uid = $6::text
                    AND  inventory.identity_library.name = $7::text)
INSERT INTO       asset.server_parent (
                         serverID,
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

	ServerTxStackClamp = `
UPDATE            asset.server_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.server_parent.serverID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.server_parent.validity;`

	ServerTxUnstackChildren = `
UPDATE            asset.runtime_environment_parent
   SET            validity = tstzrange(lower(validity), $1::timestamptz(3), '[)')
WHERE             asset.runtime_environment_parent.parentServerID = $2::uuid
  AND             $1::timestamptz(3) <@ asset.runtime_environment_parent.validity;`

	ServerTxUnstackCldClean = `
DELETE FROM       asset.runtime_environment_parent
WHERE             asset.runtime_environment_parent.parentServerID = $2::uuid
  AND             lower(asset.runtime_environment_parent.validity) > $1::timestamptz(3);`

	ServerDelNamespaceStdValues = `
DELETE FROM       asset.server_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	ServerDelNamespaceUniqValues = `
DELETE FROM       asset.server_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	ServerDelNamespace = `
DELETE FROM       asset.server
WHERE             dictionaryID = $1::uuid;`

	ServerDelNamespaceLinking = `
DELETE FROM       asset.server_linking
WHERE             dictionaryID_A = $1::uuid
   OR             dictionaryID_B = $1::uuid;`

	ServerDelNamespaceParent = `
DELETE FROM       asset.server_parent
USING             asset.server
WHERE             asset.server_parent.serverID = asset.server.serverID
  AND             asset.server.dictionaryID = $1::uuid;`
)

func init() {
	m[ServerAdd] = `ServerAdd`
	m[ServerDelNamespaceLinking] = `ServerDelNamespaceLinking`
	m[ServerDelNamespaceParent] = `ServerDelNamespaceParent`
	m[ServerDelNamespaceStdValues] = `ServerDelNamespaceStdValues`
	m[ServerDelNamespaceUniqValues] = `ServerDelNamespaceUniqValues`
	m[ServerDelNamespace] = `ServerDelNamespace`
	m[ServerLink] = `ServerLink`
	m[ServerTxStackAdd] = `ServerTxStackAdd`
	m[ServerTxStackClamp] = `ServerTxStackClamp`
	m[ServerTxStdPropertyAdd] = `ServerTxStdPropertyAdd`
	m[ServerTxStdPropertyClamp] = `ServerTxStdPropertyClamp`
	m[ServerTxStdPropertyClean] = `ServerTxStdPropertyClean`
	m[ServerTxStdPropertySelect] = `ServerTxStdPropertySelect`
	m[ServerTxUniqPropertyAdd] = `ServerTxUniqPropertyAdd`
	m[ServerTxUniqPropertyClamp] = `ServerTxUniqPropertyClamp`
	m[ServerTxUniqPropertyClean] = `ServerTxUniqPropertyClean`
	m[ServerTxUniqPropertySelect] = `ServerTxUniqPropertySelect`
	m[ServerTxUnstackChildren] = `ServerTxUnstackChildren`
	m[ServerTxUnstackCldClean] = `ServerTxUnstackCldClean`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
