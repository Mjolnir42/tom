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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
