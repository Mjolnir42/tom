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
WITH ins_dct AS ( INSERT INTO meta.dictionary ( name )
                  VALUES      ( $1::text )
                  ON CONFLICT ( name ) DO UPDATE SET name=EXCLUDED.name
                  RETURNING   dictionaryID AS dictID, name AS dictName),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute )
                  SELECT      dictID,
                              'dict_name'
                  FROM        ins_dct ),
     ins_nam AS ( INSERT INTO meta.unique_attribute ( dictionaryID, attribute )
                  SELECT      dictID,
                              'dict_name'
                  FROM        ins_dct
                  ON CONFLICT ON CONSTRAINT __uniq_unique_attr DO UPDATE SET dictionaryID=EXCLUDED.dictionaryID
                  RETURNING   attributeID AS attrID )
INSERT INTO       meta.dictionary_unique_attribute_values ( dictionaryID, attributeID, value, validity )
SELECT            dictID,
                  attrID,
                  dictName,
                  '[-infinity,infinity]'::tstzrange
FROM              ins_dct
  CROSS JOIN      ins_nam
ON CONFLICT       ON CONSTRAINT __mdq_temporal DO NOTHING;`

	NamespaceConfigure = `
WITH cte     AS ( SELECT      dictionaryID AS dictID
                  FROM        meta.dictionary
                  WHERE       name = $1::text ),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute )
                  SELECT      dictID,
                              $2::text
                  FROM        cte ),
     ins_typ AS ( INSERT INTO meta.standard_attribute ( dictionaryID, attribute )
                  SELECT      dictID,
                              $2::text
                  FROM        cte
                  ON CONFLICT ON CONSTRAINT __uniq_attribute DO UPDATE SET dictionaryID=EXCLUDED.dictionaryID
                  RETURNING   attributeID AS attrID )
INSERT INTO       meta.dictionary_standard_attribute_values ( dictionaryID, attributeID, value, validity )
SELECT            dictID,
                  attrID,
                  $3::text,
                  '[-infinity,infinity]'::tstzrange
FROM              cte
  CROSS JOIN      ins_typ
ON CONFLICT       ON CONSTRAINT __mda_temporal DO NOTHING;`

	NamespaceRemove = `
SELECT      'Namespace.REMOVE(TODO)';`

	NamespaceAttributeAddStandard = `
WITH cte     AS ( SELECT      dictionaryID AS dictID
                  FROM        meta.dictionary
                  WHERE       name = $1::text ),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute )
                  SELECT      dictID,
                              $2::text
                  FROM        cte )
INSERT INTO        meta.standard_attribute ( dictionaryID, attribute )
SELECT            cte.dictID,
                  v.value
FROM              cte
  CROSS JOIN      (VALUES($2::text)) AS v (value)
ON CONFLICT       ON CONSTRAINT __uniq_attribute DO NOTHING;`

	NamespaceAttributeAddUnique = `
WITH cte     AS ( SELECT      dictionaryID AS dictID
                  FROM        meta.dictionary
                  WHERE       name = $1::text ),
     ins_reg AS ( INSERT INTO meta.attribute ( dictionaryID, attribute )
                  SELECT      dictID,
                              $2::text
                  FROM        cte )
INSERT INTO        meta.unique_attribute ( dictionaryID, attribute )
SELECT            cte.dictID,
                  v.value
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

	NamespaceTxStdPropertySelect = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM        meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE        meta.standard_attribute.attribute = $2::text )
SELECT            value,
                  lower(validity),
                  upper(validity)
FROM              meta.dictionary_standard_attribute_values
    JOIN           cte_dct
      ON           cte_dct.dictionaryID = meta.dictionary_standard_attribute_values.dictionaryID
    JOIN          cte_att
      ON           cte_att.attributeID = meta.dictionary_standard_attribute_values.attributeID
   WHERE           now()::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity
   FOR UPDATE;`

	NamespaceTxStdPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM        meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE        meta.standard_attribute.attribute = $2::text )
UPDATE            meta.dictionary_standard_attribute_values
   SET            validity = tstzrange(lower(validity), now()::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE              meta.dictionary_standard_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND              meta.dictionary_standard_attribute_values.attributeID     = cte_att.attributeID
  AND              meta.dictionary_standard_attribute_values.value            = $3::text
  AND              lower(meta.dictionary_standard_attribute_values.validity) = $4::timestamptz(3)
  AND              upper(meta.dictionary_standard_attribute_values.validity) = $5::timestamptz(3)
  AND             now()::timestamptz(3) <@ meta.dictionary_standard_attribute_values.validity
RETURNING          upper(meta.dictionary_standard_attribute_values.validity);`

	NamespaceTxStdPropertyAdd = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.standard_attribute.attributeID
                   FROM        meta.standard_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.standard_attribute.dictionaryID
                  WHERE        meta.standard_attribute.attribute = $2::text )
INSERT INTO        meta.dictionary_standard_attribute_values ( dictionaryID, attributeID, value, validity )
SELECT            cte_dct.dictionaryID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), 'infinity', '[]')
FROM              cte_dct
    CROSS JOIN    cte_att;`

	NamespaceTxUniqPropertySelect = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM        meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE        meta.unique_attribute.attribute = $2::text )
SELECT            value,
                  lower(validity),
                  upper(validity)
FROM              meta.dictionary_unique_attribute_values
    JOIN           cte_dct
      ON           cte_dct.dictionaryID = meta.dictionary_unique_attribute_values.dictionaryID
    JOIN          cte_att
      ON           cte_att.attributeID = meta.dictionary_unique_attribute_values.attributeID
   WHERE           now()::timestamptz(3) <@ meta.dictionary_unique_attribute_values.validity
   FOR UPDATE;`

	NamespaceTxUniqPropertyClamp = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                   FROM        meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE        meta.unique_attribute.attribute = $2::text )
UPDATE            meta.dictionary_unique_attribute_values
   SET            validity = tstzrange(lower(validity), now()::timestamptz(3), '[)')
FROM              cte_dct
    CROSS JOIN    cte_att
WHERE              meta.dictionary_unique_attribute_values.dictionaryID    = cte_dct.dictionaryID
  AND              meta.dictionary_unique_attribute_values.attributeID     = cte_att.attributeID
  AND              meta.dictionary_unique_attribute_values.value            = $3::text
  AND              lower(meta.dictionary_unique_attribute_values.validity) = $4::timestamptz(3)
  AND              upper(meta.dictionary_unique_attribute_values.validity) = $5::timestamptz(3)
  AND             now()::timestamptz(3) <@ meta.dictionary_unique_attribute_values.validity
RETURNING          upper(meta.dictionary_unique_attribute_values.validity);`

	NamespaceTxUniqPropertyAdd = `
WITH cte_dct AS ( SELECT      meta.dictionary.dictionaryID
                  FROM        meta.dictionary
                  WHERE       meta.dictionary.name = $1::text ),
     cte_att AS ( SELECT      meta.unique_attribute.attributeID
                  FROM        meta.unique_attribute
                    JOIN      cte_dct
                      ON      cte_dct.dictionaryID = meta.unique_attribute.dictionaryID
                  WHERE        meta.unique_attribute.attribute = $2::text )
INSERT INTO       meta.dictionary_unique_attribute_values ( dictionaryID, attributeID, value, validity )
SELECT            cte_dct.dictionaryID,
                  cte_att.attributeID,
                  $3::text,
                  tstzrange($4::timestamptz(3), 'infinity', '[]')
FROM              cte_dct
    CROSS JOIN    cte_att;`
)

func init() {
	m[NamespaceAdd] = `NamespaceAdd`
	m[NamespaceAttributeAddStandard] = `NamespaceAttributeAddStandard`
	m[NamespaceAttributeAddUnique] = `NamespaceAttributeAddUnique`
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
