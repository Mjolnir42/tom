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

func init() {
	Commands[CmdLibraryAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/idlib`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdLibraryList] = CmdDef{
		Method:      MethodGET,
		Path:        `/idlib`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdLibraryRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/idlib/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdLibraryShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/idlib/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
}

// Library ...
type Library struct {
	ID        string `json:"-"`
	Name      string `json:"name"`
	IsSelfEnr bool   `json:"is-selfenrollment-enabled"`
	IsMachine bool   `json:"is-machine-library"`
	EnrolKey  string `json:"enrolment-key"`
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
	TomID     string `json:"-"`
}

func (l *Library) SetTomID() Entity {
	l.TomID = l.FormatDNS()
	return l
}

func (l *Library) String() string {
	return l.FormatDNS()
}

func (l *Library) FormatDNS() string {
	return l.Name + `.` + EntityLibrary + `.tom`
}

func (l *Library) FormatTomID() string {
	return `tom:///` + EntityLibrary + `/name=` + l.Name
}

func (l *Library) ParseTomID() error {
	var typeID string
	switch {
	case l.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(l.TomID):
		l.Name, _, typeID = parseTomIDFormatDNS(l.TomID)
		return assessTomID(EntityLibrary, typeID)
	case isTomIDFormatURI(l.TomID):
		l.Name, _, typeID = parseTomIDFormatURI(l.TomID)
		return assessTomID(EntityLibrary, typeID)
	default:
		return ErrInvalidTomID
	}
}

func (l *Library) ExportName() string {
	return l.Name
}

func (l *Library) ExportNamespace() string {
	return l.Name
}

func (l *Library) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		close(ret)
	}()
	return ret
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
