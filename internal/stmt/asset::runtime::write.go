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
  CROSS JOIN      sel_dct;`

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

	RuntimeTxStdPropertyAdd     = `SELECT 'RuntimeTxStdPropertyAdd';`
	RuntimeTxStdPropertyClamp   = `SELECT 'RuntimeTxStdPropertyClamp';`
	RuntimeTxStdPropertySelect  = `SELECT 'RuntimeTxStdPropertySelect';`
	RuntimeTxUniqPropertyAdd    = `SELECT 'RuntimeTxUniqPropertyAdd';`
	RuntimeTxUniqPropertyClamp  = `SELECT 'RuntimeTxUniqPropertyClamp';`
	RuntimeTxUniqPropertySelect = `SELECT 'RuntimeTxUniqPropertySelect';`
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
