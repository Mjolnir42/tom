/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Socket ...
type Socket struct {
	Namespace    string                    `json:"namespace"`
	Name         string                    `json:"name"`
	Parent       string                    `json:"parent"`
	Property     map[string]PropertyDetail `json:"property"`
	CreatedAt    string                    `json:"createdAt"`
	CreatedBy    string                    `json:"createdBy"`
	ID           string                    `json:"-"`
	TomID        string                    `json:"-"`
	StdProperty  []PropertyDetail          `json:"-"`
	UniqProperty []PropertyDetail          `json:"-"`
}

func (s *Socket) SetTomID() Entity {
	s.TomID = s.FormatDNS()
	return s
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

func (s *Socket) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		for key := range s.Property {
			ret <- s.Property[key]
		}
		close(ret)
	}()
	return ret
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
