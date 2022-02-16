/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdLibrary       = ModelIAM + `::` + EntityLibrary + `:`
	CmdLibraryAdd    = ModelIAM + `::` + EntityLibrary + `:` + ActionAdd
	CmdLibraryList   = ModelIAM + `::` + EntityLibrary + `:` + ActionList
	CmdLibraryRemove = ModelIAM + `::` + EntityLibrary + `:` + ActionRemove
	CmdLibraryShow   = ModelIAM + `::` + EntityLibrary + `:` + ActionShow
)

// Library ...
type Library struct {
	ID        string `json:"-"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
