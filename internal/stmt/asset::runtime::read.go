/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	RuntimeReadStatements = ``

	RuntimeList = `
SELECT            meta.dictionary.name AS dictionaryName,
                  asset.runtime_environment_unique_attribute_values.value AS runtimeName,
                  inventory.user.uid AS createdBy,
                  asset.runtime_environment_unique_attribute_values.createdAt
FROM              meta.dictionary
JOIN              meta.unique_attribute
  ON              meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
JOIN              asset.runtime_environment
  ON              meta.dictionary.dictionaryID = asset.runtime_environment.dictionaryID
JOIN              asset.runtime_environment_unique_attribute_values
    ON            meta.dictionary.dictionaryID = asset.runtime_environment_unique_attribute_values.dictionaryID
    AND           asset.runtime_environment.rteID = asset.runtime_environment_unique_attribute_values.rteID
    AND           meta.unique_attribute.attributeID = asset.runtime_environment_unique_attribute_values.attributeID
JOIN              inventory.user
  ON              asset.runtime_environment_unique_attribute_values.createdBy = inventory.user.userID
WHERE             (meta.dictionary.name = $1::text OR $1::text IS NULL)
  AND             meta.unique_attribute.attribute = 'name'::text
  AND             now()::timestamptz(3) <@ asset.runtime_environment_unique_attribute_values.validity;`

	RuntimeShow = `
SELECT      'Runtime.SHOW';`
)

func init() {
	m[RuntimeList] = `RuntimeList`
	m[RuntimeShow] = `RuntimeShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
