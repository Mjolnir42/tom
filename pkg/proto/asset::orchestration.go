/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdOrchestration           = ModelAsset + `::` + EntityOrchestration + `:`
	CmdOrchestrationAdd        = ModelAsset + `::` + EntityOrchestration + `:` + ActionAdd
	CmdOrchestrationLink       = ModelAsset + `::` + EntityOrchestration + `:` + ActionLink
	CmdOrchestrationList       = ModelAsset + `::` + EntityOrchestration + `:` + ActionList
	CmdOrchestrationPropRemove = ModelAsset + `::` + EntityOrchestration + `:` + ActionPropRemove
	CmdOrchestrationPropSet    = ModelAsset + `::` + EntityOrchestration + `:` + ActionPropSet
	CmdOrchestrationPropUpdate = ModelAsset + `::` + EntityOrchestration + `:` + ActionPropUpdate
	CmdOrchestrationRemove     = ModelAsset + `::` + EntityOrchestration + `:` + ActionRemove
	CmdOrchestrationResolve    = ModelAsset + `::` + EntityOrchestration + `:` + ActionResolve
	CmdOrchestrationShow       = ModelAsset + `::` + EntityOrchestration + `:` + ActionShow
	CmdOrchestrationStack      = ModelAsset + `::` + EntityOrchestration + `:` + ActionStack
	CmdOrchestrationUnstack    = ModelAsset + `::` + EntityOrchestration + `:` + ActionUnstack
)

func init() {
	Commands[CmdOrchestrationAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/orchestration/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdOrchestrationList] = CmdDef{
		Method:      MethodGET,
		Path:        `/orchestration/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdOrchestrationLink] = CmdDef{
		Method:      MethodPOST,
		Path:        `/orchestration/` + PlHoldTomID + `/link/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationPropSet] = CmdDef{
		Method:      MethodPUT,
		Path:        `/orchestration/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationPropUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/orchestration/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationPropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/orchestration/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/orchestration/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/orchestration/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationStack] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/orchestration/` + PlHoldTomID + `/parent`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationUnstack] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/orchestration/` + PlHoldTomID + `/parent`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdOrchestrationResolve] = CmdDef{
		Method:      MethodGET,
		Path:        `/orchestration/` + PlHoldTomID + `/resolve/` + PlHoldResolv,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{PlHoldTomID, PlHoldResolv},
	}
}

// Orchestration defines a orchestration environment within the asset model
type Orchestration struct {
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Type      string                    `json:"type"`
	Parent    []string                  `json:"parent"`
	Link      []string                  `json:"link"`
	Children  []string                  `json:"children"`
	Property  map[string]PropertyDetail `json:"property"`
	CreatedAt string                    `json:"createdAt"`
	CreatedBy string                    `json:"createdBy"`
	ID        string                    `json:"-"`
	TomID     string                    `json:"-"`
}

func NewOrchestrationRequest() Request {
	return Request{
		Orchestration: NewOrchestration(),
	}
}

func NewOrchestration() *Orchestration {
	return &Orchestration{
		Parent:   []string{},
		Link:     []string{},
		Children: []string{},
		Property: map[string]PropertyDetail{},
	}
}

// OrchestrationHeader defines ....
type OrchestrationHeader struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
}

func (o *Orchestration) SetTomID() Entity {
	o.TomID = o.FormatDNS()
	return o
}

func (o *Orchestration) String() string {
	return o.FormatDNS()
}

func (o *Orchestration) FormatDNS() string {
	return o.Name + `.` + o.Namespace + `.` + EntityOrchestration + `.tom`
}

func (o *Orchestration) FormatTomID() string {
	return `tom://` + o.Namespace + `/` + EntityOrchestration + `/name=` + o.Name
}

func (o *Orchestration) ParseTomID() error {
	var typeID string
	switch {
	case o.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(o.TomID):
		o.Name, o.Namespace, typeID = parseTomIDFormatDNS(o.TomID)
		return assessTomID(EntityOrchestration, typeID)
	case isTomIDFormatURI(o.TomID):
		o.Name, o.Namespace, typeID = parseTomIDFormatURI(o.TomID)
		return assessTomID(EntityOrchestration, typeID)
	default:
		return ErrInvalidTomID
	}
}

func (o *Orchestration) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		for key := range o.Property {
			ret <- o.Property[key]
		}
		close(ret)
	}()
	return ret
}

func (o *Orchestration) ExportName() string {
	return o.Name
}

func (o *Orchestration) ExportNamespace() string {
	return o.Namespace
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
