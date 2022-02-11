/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// User ...
type Team struct {
	LibraryName string  `json:"library-name"`
	TeamName    string  `json:"team-name"`
	ExternalID  string  `json:"external-ref,omitempty"`
	IsDeleted   bool    `json:"is-deleted"`
	CreatedAt   string  `json:"createdAt"`
	CreatedBy   string  `json:"createdBy"`
	TeamLead    *User   `json:"team-lead,omitempty"`
	Member      *[]User `json:"team-member,omitempty"`
	ID          string  `json:"-"`
	LibraryID   string  `json:"-"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
