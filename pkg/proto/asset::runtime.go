/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdRuntime           = ModelAsset + `::` + EntityRuntime + `:`
	CmdRuntimeAdd        = ModelAsset + `::` + EntityRuntime + `:` + ActionAdd
	CmdRuntimeLink       = ModelAsset + `::` + EntityRuntime + `:` + ActionLink
	CmdRuntimeList       = ModelAsset + `::` + EntityRuntime + `:` + ActionList
	CmdRuntimePropRemove = ModelAsset + `::` + EntityRuntime + `:` + ActionPropRemove
	CmdRuntimePropSet    = ModelAsset + `::` + EntityRuntime + `:` + ActionPropSet
	CmdRuntimePropUpdate = ModelAsset + `::` + EntityRuntime + `:` + ActionPropUpdate
	CmdRuntimeRemove     = ModelAsset + `::` + EntityRuntime + `:` + ActionRemove
	CmdRuntimeResolve    = ModelAsset + `::` + EntityRuntime + `:` + ActionResolve
	CmdRuntimeShow       = ModelAsset + `::` + EntityRuntime + `:` + ActionShow
	CmdRuntimeStack      = ModelAsset + `::` + EntityRuntime + `:` + ActionStack
	CmdRuntimeUnstack    = ModelAsset + `::` + EntityRuntime + `:` + ActionUnstack
)

func init() {
	Commands[CmdRuntimeAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/runtime/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdRuntimeLink] = CmdDef{
		Method:      MethodPOST,
		Path:        `/runtime/` + PlHoldTomID + `/link/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimeList] = CmdDef{
		Method:      MethodGET,
		Path:        `/runtime/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdRuntimePropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/runtime/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimePropSet] = CmdDef{
		Method:      MethodPUT,
		Path:        `/runtime/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimePropUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/runtime/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimeRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/runtime/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimeResolve] = CmdDef{
		Method:      MethodGET,
		Path:        `/runtime/` + PlHoldTomID + `/resolve/` + PlHoldResolv,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{PlHoldTomID, PlHoldResolv},
	}
	Commands[CmdRuntimeShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/runtime/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimeStack] = CmdDef{
		Method:      MethodPUT,
		Path:        `/runtime/` + PlHoldTomID + `/parent`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdRuntimeUnstack] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/runtime/` + PlHoldTomID + `/parent`,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
}

// Runtime defines a runtime within the asset model
type Runtime struct {
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

func NewRuntimeRequest() Request {
	return Request{
		Runtime: NewRuntime(),
	}
}

func NewRuntime() *Runtime {
	return &Runtime{
		Link:      []string{},
		Children:  []string{},
		Property:  map[string]PropertyDetail{},
		Resources: []string{},
	}
}

// RuntimeHeader defines ....
type RuntimeHeader struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

func (r *Runtime) SetTomID() Entity {
	r.TomID = r.FormatDNS()
	return r
}

func (r *Runtime) String() string {
	return r.FormatDNS()
}

func (r *Runtime) FormatDNS() string {
	return r.Name + `.` + r.Namespace + `.` + EntityRuntime + `.tom`
}

func (r *Runtime) FormatTomID() string {
	return `tom://` + r.Namespace + `/` + EntityRuntime + `/name=` + r.Name
}

func (r *Runtime) ParseTomID() error {
	var typeID string
	switch {
	case r.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(r.TomID):
		r.Name, r.Namespace, typeID = parseTomIDFormatDNS(r.TomID)
		return assessTomID(EntityRuntime, typeID)
	case isTomIDFormatURI(r.TomID):
		r.Name, r.Namespace, typeID = parseTomIDFormatURI(r.TomID)
		return assessTomID(EntityRuntime, typeID)
	default:
		return ErrInvalidTomID
	}
}

func (r *Runtime) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		for key := range r.Property {
			ret <- r.Property[key]
		}
		close(ret)
	}()
	return ret
}

func (r *Runtime) ExportName() string {
	return r.Name
}

func (r *Runtime) ExportNamespace() string {
	return r.Namespace
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
