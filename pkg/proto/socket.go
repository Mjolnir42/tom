/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const EntitySocket = `socket`

// Socket ...
type Socket struct {
	ID           string            `json:"-"`
	TomID        string            `json:"-"`
	Namespace    string            `json:"namespace"`
	Name         string            `json:"name"`
	Parent       string            `json:"parent"`
	PropertyMap  map[string]string `json:"property"`
	StdProperty  []Property        `json:"-"`
	UniqProperty []Property        `json:"-"`
}

func (s *Socket) String() string {
	return s.FormatDNS()
}

func (s *Socket) FormatDNS() string {
	return s.Name + `.` + s.Namespace + `.` + EntityContainer + `.tom`
}

func (s *Socket) FormatTomID() string {
	return `tom://` + s.Namespace + `/` + EntityContainer + `/name=` + s.Name
}

func (s *Socket) ParseTomID() error {
	switch {
	case s.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(s.TomID):
		s.Name, s.Namespace, _ = parseTomIDFormatDNS(s.TomID)
		return nil
	case isTomIDFormatURI(s.TomID):
		s.Name, s.Namespace, _ = parseTomIDFormatURI(s.TomID)
		return nil
	default:
		return ErrInvalidTomID
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
