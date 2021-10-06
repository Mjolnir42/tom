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
SELECT      'Server.ADD';`

	ServerRemove = `
SELECT      'Server.REMOVE';`

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

	ServerStdAttrRemove = `
DELETE FROM       asset.server_standard_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

	ServerUniqAttrRemove = `
DELETE FROM       asset.server_unique_attribute_values
WHERE             attributeID = $1::uuid
  AND             dictionaryID = $2::uuid;`

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
)

func init() {
	m[ServerAdd] = `ServerAdd`
	m[ServerLink] = `ServerLink`
	m[ServerRemove] = `ServerRemove`
	m[ServerStdAttrRemove] = `ServerStdAttrRemove`
	m[ServerTxStdPropertyAdd] = `ServerTxStdPropertyAdd`
	m[ServerTxStdPropertyClamp] = `ServerTxStdPropertyClamp`
	m[ServerTxStdPropertySelect] = `ServerTxStdPropertySelect`
	m[ServerTxUniqPropertyAdd] = `ServerTxUniqPropertyAdd`
	m[ServerTxUniqPropertyClamp] = `ServerTxUniqPropertyClamp`
	m[ServerTxUniqPropertySelect] = `ServerTxUniqPropertySelect`
	m[ServerUniqAttrRemove] = `ServerUniqAttrRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
