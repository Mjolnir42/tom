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
                  FROM        cte ),
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

	NamespaceRemove = `
SELECT      'Namespace.REMOVE';`
)

func init() {
	m[NamespaceAdd] = `NamespaceAdd`
	m[NamespaceRemove] = `NamespaceRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
