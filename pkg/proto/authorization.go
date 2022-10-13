/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Request is the request wrapper of Tom's public API
type Authorization struct {
	Timestamp string `json:"timestamp"`
	UserID    string `json:"userID"`
	DataHash  string `json:"dataHash"`
	Signature string `json:"signature"`
}

// Serialize ...
func (a *Authorization) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(a.Timestamp)...)
	data = append(data, []byte(a.UserID)...)
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
