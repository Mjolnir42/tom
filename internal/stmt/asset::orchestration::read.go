/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	OrchestrationReadStatements = ``

	OrchestrationTxShow = `
SELECT            asset.orchestration_environment.orchID,
                  asset.orchestration_environment.dictionaryID,
                  asset.orchestration_environment.createdAt,
                  creator.uid AS createdBy,
                  lower(asset.orchestration_environment_unique_attribute_values.validity) AS validSince,
                  upper(asset.orchestration_environment_unique_attribute_values.validity) AS validUntil,
                  asset.orchestration_environment_unique_attribute_values.createdAt AS namedAt,
                  namegiver.uid AS namedBy
FROM              meta.dictionary
    JOIN          asset.orchestration_environment
        ON        meta.dictionary.dictionaryID = asset.orchestration_environment.dictionaryID
    JOIN          inventory.user AS creator
        ON        asset.orchestration_environment.createdBy = creator.userID
    JOIN          meta.unique_attribute
        ON        meta.dictionary.dictionaryID = meta.unique_attribute.dictionaryID
    JOIN          asset.orchestration_environment_unique_attribute_values
        ON        meta.dictionary.dictionaryID = asset.orchestration_environment_unique_attribute_values.dictionaryID
        AND       asset.orchestration_environment.orchID = asset.orchestration_environment_unique_attribute_values.orchID
        AND       meta.unique_attribute.attributeID = asset.orchestration_environment_unique_attribute_values.attributeID
    JOIN          inventory.user AS namegiver
        ON        asset.orchestration_environment_unique_attribute_values.createdBy = namegiver.userID
WHERE             meta.dictionary.name = $1::text
     AND          meta.unique_attribute.attribute = 'name'::text
     AND          asset.orchestration_environment_unique_attribute_values.value = $2::text
     AND          $3::timestamptz(3) <@ asset.orchestration_environment_unique_attribute_values.validity;`
)

func init() {
	m[OrchestrationTxShow] = `OrchestrationTxShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
