/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdContainer           = ModelAsset + `::` + EntityContainer + `:`
	CmdContainerAdd        = ModelAsset + `::` + EntityContainer + `:` + ActionAdd
	CmdContainerList       = ModelAsset + `::` + EntityContainer + `:` + ActionList
	CmdContainerLink       = ModelAsset + `::` + EntityContainer + `:` + ActionLink
	CmdContainerPropRemove = ModelAsset + `::` + EntityContainer + `:` + ActionPropRemove
	CmdContainerPropSet    = ModelAsset + `::` + EntityContainer + `:` + ActionPropSet
	CmdContainerPropUpdate = ModelAsset + `::` + EntityContainer + `:` + ActionPropUpdate
	CmdContainerRemove     = ModelAsset + `::` + EntityContainer + `:` + ActionRemove
	CmdContainerShow       = ModelAsset + `::` + EntityContainer + `:` + ActionShow
)

func init() {
	Commands[CmdContainerAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/container/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdContainerList] = CmdDef{
		Method:      MethodGET,
		Path:        `/container/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdContainerLink] = CmdDef{
		Method:      MethodPOST,
		Path:        `/container/` + PlHoldTomID + `/link/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerPropSet] = CmdDef{
		Method:      MethodPUT,
		Path:        `/container/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerPropUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/container/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerPropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/container/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/container/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/container/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
}

// Container ...
type Container struct {
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

func NewContainerRequest() Request {
	return Request{
		Container: &Container{
			Link:     []string{},
			Property: map[string]PropertyDetail{},
		},
	}
}

type ContainerHeader struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
}

func (c *Container) SetTomID() Entity {
	c.TomID = c.FormatDNS()
	return c
}

func (c *Container) String() string {
	return c.FormatDNS()
}

func (c *Container) FormatDNS() string {
	return c.Name + `.` + c.Namespace + `.` + EntityContainer + `.tom`
}

func (c *Container) FormatTomID() string {
	return `tom://` + c.Namespace + `/` + EntityContainer + `/name=` + c.Name
}

func (c *Container) ParseTomID() error {
	var typeID string
	switch {
	case c.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(c.TomID):
		c.Name, c.Namespace, typeID = parseTomIDFormatDNS(c.TomID)
		return assessTomID(EntityContainer, typeID)
	case isTomIDFormatURI(c.TomID):
		c.Name, c.Namespace, typeID = parseTomIDFormatURI(c.TomID)
		return assessTomID(EntityContainer, typeID)
	default:
		return ErrInvalidTomID
	}
}

func (c *Container) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		for key := range c.Property {
			ret <- c.Property[key]
		}
		close(ret)
	}()
	return ret
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
