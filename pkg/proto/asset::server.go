/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdServer           = ModelAsset + `::` + EntityServer + `:`
	CmdServerAdd        = ModelAsset + `::` + EntityServer + `:` + ActionAdd
	CmdServerLink       = ModelAsset + `::` + EntityServer + `:` + ActionLink
	CmdServerList       = ModelAsset + `::` + EntityServer + `:` + ActionList
	CmdServerPropRemove = ModelAsset + `::` + EntityServer + `:` + ActionPropRemove
	CmdServerPropSet    = ModelAsset + `::` + EntityServer + `:` + ActionPropSet
	CmdServerPropUpdate = ModelAsset + `::` + EntityServer + `:` + ActionPropUpdate
	CmdServerRemove     = ModelAsset + `::` + EntityServer + `:` + ActionRemove
	CmdServerResolve    = ModelAsset + `::` + EntityServer + `:` + ActionResolve
	CmdServerShow       = ModelAsset + `::` + EntityServer + `:` + ActionShow
	CmdServerStack      = ModelAsset + `::` + EntityServer + `:` + ActionStack
	CmdServerUnstack    = ModelAsset + `::` + EntityServer + `:` + ActionUnstack
)

func init() {
	Commands[CmdServerAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/server/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdServerLink] = CmdDef{
		Method:      MethodPOST,
		Path:        `/server/` + PlHoldTomID + `/link/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerList] = CmdDef{
		Method:      MethodGET,
		Path:        `/server/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
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
	Commands[CmdServerPropUpdate] = CmdDef{
		Method:      MethodPATCH,
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
	Commands[CmdServerResolve] = CmdDef{
		Method:      MethodGET,
		Path:        `/server/` + PlHoldTomID + `/resolve/` + PlHoldResolv,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{PlHoldTomID, PlHoldResolv},
	}
	Commands[CmdServerShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/server/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerStack] = CmdDef{
		Method:      MethodPUT,
		Path:        `/server/` + PlHoldTomID + `/parent`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdServerUnstack] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/server/` + PlHoldTomID + `/parent`,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
}

// Server defines a server within the asset model
type Server struct {
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Type      string                    `json:"type"`
	Parent    string                    `json:"parent"`
	Link      []string                  `json:"link"`
	Children  []string                  `json:"children"`
	Property  map[string]PropertyDetail `json:"property"`
	CreatedAt string                    `json:"createdAt,omitempty"`
	CreatedBy string                    `json:"createdBy,omitempty"`
	Resources []string                  `json:"resources"`
	ID        string                    `json:"-"`
	TomID     string                    `json:"-"`
}

func NewServerRequest() Request {
	return Request{
		Server: NewServer(),
	}
}

func NewServer() *Server {
	return &Server{
		Link:      []string{},
		Children:  []string{},
		Property:  map[string]PropertyDetail{},
		Resources: []string{},
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

func (s *Server) ExportName() string {
	return s.Name
}

func (s *Server) ExportNamespace() string {
	return s.Namespace
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
