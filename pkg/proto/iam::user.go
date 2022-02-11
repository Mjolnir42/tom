/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// User ...
type User struct {
	LibraryName    string `json:"library-name"`
	FirstName      string `json:"first-name,omitempty"`
	LastName       string `json:"last-name,omitempty"`
	UserName       string `json:"user-name"`
	EmployeeNumber string `json:"employee-number,omitempty"`
	MailAddress    string `json:"mailaddress,omitempty"`
	ExternalID     string `json:"external-ref,omitempty"`
	IsActive       bool   `json:"is-active"`
	IsDeleted      bool   `json:"is-deleted"`
	CreatedAt      string `json:"createdAt"`
	CreatedBy      string `json:"createdBy"`
	ID             string `json:"-"`
	LibraryID      string `json:"-"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
