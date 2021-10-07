/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

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
}

// Orchestration defines a orchestration environment within the asset model
type Orchestration struct {
	Namespace    string                    `json:"namespace"`
	Name         string                    `json:"name"`
	Type         string                    `json:"type"`
	Parent       []string                  `json:"parent"`
	Link         []string                  `json:"link"`
	Property     map[string]PropertyDetail `json:"property"`
	CreatedAt    string                    `json:"createdAt"`
	CreatedBy    string                    `json:"createdBy"`
	ID           string                    `json:"-"`
	TomID        string                    `json:"-"`
	StdProperty  []PropertyDetail          `json:"-"`
	UniqProperty []PropertyDetail          `json:"-"`
}

// OrchestrationHeader defines ....
type OrchestrationHeader struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
