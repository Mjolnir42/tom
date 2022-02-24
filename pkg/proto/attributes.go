/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"encoding/json"
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
