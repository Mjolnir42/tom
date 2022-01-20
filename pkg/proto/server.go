/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

func init() {
	Commands[CmdServerAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/server/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdServerList] = CmdDef{
		Method:      MethodGET,
		Path:        `/server/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdServerLink] = CmdDef{
		Method:      MethodPOST,
		Path:        `/server/` + PlHoldTomID + `/link/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerPropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/server/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerPropSet] = CmdDef{
		Method:      MethodPUT,
		Path:        `/server/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/server/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerPropUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/server/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/server/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
}

// Server defines a server within the asset model
type Server struct {
	Namespace    string                    `json:"namespace"`
	Name         string                    `json:"name"`
	Type         string                    `json:"type"`
	Parent       string                    `json:"parent"`
	Link         []string                  `json:"link"`
	Property     map[string]PropertyDetail `json:"property"`
	CreatedAt    string                    `json:"createdAt"`
	CreatedBy    string                    `json:"createdBy"`
	ID           string                    `json:"-"`
	TomID        string                    `json:"-"`
	StdProperty  []PropertyDetail          `json:"-"`
	UniqProperty []PropertyDetail          `json:"-"`
}

func NewServerRequest() Request {
	return Request{
		Server: &Server{
			Link:     []string{},
			Property: map[string]PropertyDetail{},
		},
	}
}

// ServerHeader defines ....
type ServerHeader struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
}

func (s *Server) SetTomID() Entity {
	s.TomID = s.FormatDNS()
	return s
}

func (s *Server) String() string {
	return s.FormatDNS()
}

func (s *Server) FormatDNS() string {
	return s.Name + `.` + s.Namespace + `.` + EntityServer + `.tom`
}

func (s *Server) FormatTomID() string {
	return `tom://` + s.Namespace + `/` + EntityServer + `/name=` + s.Name
}

func (s *Server) ParseTomID() error {
	var typeID string
	switch {
	case s.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(s.TomID):
		s.Name, s.Namespace, typeID = parseTomIDFormatDNS(s.TomID)
		return assessTomID(EntityServer, typeID)
	case isTomIDFormatURI(s.TomID):
		s.Name, s.Namespace, typeID = parseTomIDFormatURI(s.TomID)
		return assessTomID(EntityServer, typeID)
	default:
		return ErrInvalidTomID
	}
}

func (s *Server) PropertyIterator() <-chan PropertyDetail {
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
