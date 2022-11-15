/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdFlow           = ModelBulk + `::` + EntityFlow + `:`
	CmdFlowAdd        = ModelBulk + `::` + EntityFlow + `:` + ActionAdd
	CmdFlowEnsure     = ModelBulk + `::` + EntityFlow + `:` + ActionEnsure
	CmdFlowList       = ModelBulk + `::` + EntityFlow + `:` + ActionList
	CmdFlowPropRemove = ModelBulk + `::` + EntityFlow + `:` + ActionPropRemove
	CmdFlowPropSet    = ModelBulk + `::` + EntityFlow + `:` + ActionPropSet
	CmdFlowPropUpdate = ModelBulk + `::` + EntityFlow + `:` + ActionPropUpdate
	CmdFlowRemove     = ModelBulk + `::` + EntityFlow + `:` + ActionRemove
	CmdFlowShow       = ModelBulk + `::` + EntityFlow + `:` + ActionShow
)

func init() {
	Commands[CmdFlowAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/flow/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdFlowEnsure] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/flow/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdFlowList] = CmdDef{
		Method:      MethodGET,
		Path:        `/flow/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdFlowPropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/flow/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdFlowPropSet] = CmdDef{
		Method:      MethodPUT,
		Path:        `/flow/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdFlowPropUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/flow/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdFlowRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/flow/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdFlowShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/flow/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
}

// Flow ...
type Flow struct {
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Type      string                    `json:"type"`
	Property  map[string]PropertyDetail `json:"property"`
	CreatedAt string                    `json:"createdAt,omitempty"`
	CreatedBy string                    `json:"createdBy,omitempty"`
	ID        string                    `json:"-"`
	TomID     string                    `json:"-"`
}

func NewFlowRequest() Request {
	return Request{
		Flow: NewFlow(),
	}
}

func NewFlow() *Flow {
	return &Flow{
		Property: map[string]PropertyDetail{},
	}
}

func (c *Flow) SetTomID() Entity {
	c.TomID = c.FormatDNS()
	return c
}

func (c *Flow) String() string {
	return c.FormatDNS()
}

func (c *Flow) FormatDNS() string {
	return c.Name + `.` + c.Namespace + `.` + EntityFlow + `.tom`
}

func (c *Flow) FormatTomID() string {
	return `tom://` + c.Namespace + `/` + EntityFlow + `/name=` + c.Name
}

func (c *Flow) ParseTomID() error {
	var typeID string
	switch {
	case c.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(c.TomID):
		c.Name, c.Namespace, typeID = parseTomIDFormatDNS(c.TomID)
		return assessTomID(EntityFlow, typeID)
	case isTomIDFormatURI(c.TomID):
		c.Name, c.Namespace, typeID = parseTomIDFormatURI(c.TomID)
		return assessTomID(EntityFlow, typeID)
	default:
		return ErrInvalidTomID
	}
}

func (c *Flow) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		for key := range c.Property {
			ret <- c.Property[key]
		}
		close(ret)
	}()
	return ret
}

func (c *Flow) ExportName() string {
	return c.Name
}

func (c *Flow) ExportNamespace() string {
	return c.Namespace
}

// Serialize ...
func (c *Flow) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(c.Namespace)...)
	data = append(data, []byte(c.Name)...)
	data = append(data, []byte(c.Type)...)
	data = append(data, SerializeMapPropertyDetail(c.Property)...)
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
