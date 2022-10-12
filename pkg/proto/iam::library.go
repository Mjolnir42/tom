/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
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
	IsSelfEnr bool   `json:"is-selfenrollment-enabled"`
	IsMachine bool   `json:"is-machine-library"`
	EnrollKey string `json:"-"`
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

// Serialize ...
func (l *Library) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(l.Name)...)
	switch l.IsSelfEnr {
	case true:
		data = append(data, []byte(`true`)...)
	default:
		data = append(data, []byte(`false`)...)
	}
	switch l.IsMachine {
	case true:
		data = append(data, []byte(`true`)...)
	default:
		data = append(data, []byte(`false`)...)
	}
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
