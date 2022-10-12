/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"encoding/json"
	"sort"
)

// PropertyDetail holds all the information about an object property
type PropertyDetail struct {
	Attribute  string          `json:"attribute"`
	Value      string          `json:"value"`
	Raw        json.RawMessage `json:"structuredValue,omitempty"`
	ValidSince string          `json:"validSince,omitempty"`
	ValidUntil string          `json:"validUntil,omitempty"`
	CreatedAt  string          `json:"createdAt,omitempty"`
	CreatedBy  string          `json:"createdBy,omitempty"`
	Namespace  string          `json:"namespace"`
}

// AttributeDefinition holds the definition of a dictionary attribute
type AttributeDefinition struct {
	Key    string `json:"key"`
	Unique bool   `json:"uniqueValueConstraint"`
}

// SerializeMapPropertyDetail
func SerializeMapPropertyDetail(m map[string]PropertyDetail) []byte {
	data := make([]byte, 0)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		data = append(data, []byte(m[k].Attribute)...)
		data = append(data, []byte(m[k].Value)...)
		data = append(data, []byte(m[k].Namespace)...)
		data = append(data, []byte(m[k].ValidSince)...)
		data = append(data, []byte(m[k].ValidUntil)...)
	}
	return data
}

// SerializeAttributeSlice ...
func SerializeAttributeSlice(s []AttributeDefinition) []byte {
	data := make([]byte, 0)
	for _, def := range s {
		data = append(data, []byte(def.Key)...)
		switch def.Unique {
		case true:
			data = append(data, []byte(`true`)...)
		default:
			data = append(data, []byte(`false`)...)
		}
	}
	return data
}

// SerializeStringSlice ...
func SerializeStringSlice(s []string) []byte {
	data := make([]byte, 0)
	for _, str := range s {
		data = append(data, []byte(str)...)
	}
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
