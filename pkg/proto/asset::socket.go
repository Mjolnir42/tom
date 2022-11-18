/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Socket ...
type Socket struct {
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Parent    string                    `json:"parent"`
	Property  map[string]PropertyDetail `json:"property"`
	CreatedAt string                    `json:"createdAt,omitempty"`
	CreatedBy string                    `json:"createdBy,omitempty"`
	ID        string                    `json:"-"`
	TomID     string                    `json:"-"`
}

func (s *Socket) SetTomID() Entity {
	s.TomID = s.FormatDNS()
	return s
}

func (s *Socket) String() string {
	return s.FormatDNS()
}

func (s *Socket) FormatDNS() string {
	return s.Name + `.` + s.Namespace + `.` + EntitySocket + `.tom`
}

func (s *Socket) FormatTomID() string {
	return `tom://` + s.Namespace + `/` + EntitySocket + `/name=` + s.Name
}

func (s *Socket) ParseTomID() error {
	var typeID string
	switch {
	case s.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(s.TomID):
		s.Name, s.Namespace, typeID = parseTomIDFormatDNS(s.TomID)
		return assessTomID(EntitySocket, typeID)
	case isTomIDFormatURI(s.TomID):
		s.Name, s.Namespace, typeID = parseTomIDFormatURI(s.TomID)
		return assessTomID(EntitySocket, typeID)
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

func (s *Socket) ExportName() string {
	return s.Name
}

func (s *Socket) ExportNamespace() string {
	return s.Namespace
}

// Serialize ...
func (s *Socket) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(s.Namespace)...)
	data = append(data, []byte(s.Name)...)
	data = append(data, []byte(s.Parent)...)
	data = append(data, SerializeMapPropertyDetail(s.Property)...)
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
