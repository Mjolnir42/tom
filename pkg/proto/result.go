/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Result is the response wrapper of Tom's public API
type Result struct {
	Namespace     *[]Namespace     `json:"namespace,omitempty"`
	Orchestration *[]Orchestration `json:"orchestration,omitempty"`
	Runtime       *[]Runtime       `json:"runtime,omitempty"`
	Server        *[]Server        `json:"server,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
