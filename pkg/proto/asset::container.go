/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdContainer           = ModelAsset + `::` + EntityContainer + `:`
	CmdContainerAdd        = ModelAsset + `::` + EntityContainer + `:` + ActionAdd
	CmdContainerLink       = ModelAsset + `::` + EntityContainer + `:` + ActionLink
	CmdContainerList       = ModelAsset + `::` + EntityContainer + `:` + ActionList
	CmdContainerPropRemove = ModelAsset + `::` + EntityContainer + `:` + ActionPropRemove
	CmdContainerPropSet    = ModelAsset + `::` + EntityContainer + `:` + ActionPropSet
	CmdContainerPropUpdate = ModelAsset + `::` + EntityContainer + `:` + ActionPropUpdate
	CmdContainerRemove     = ModelAsset + `::` + EntityContainer + `:` + ActionRemove
	CmdContainerResolve    = ModelAsset + `::` + EntityContainer + `:` + ActionResolve
	CmdContainerShow       = ModelAsset + `::` + EntityContainer + `:` + ActionShow
	CmdContainerStack      = ModelAsset + `::` + EntityContainer + `:` + ActionStack
	CmdContainerUnstack    = ModelAsset + `::` + EntityContainer + `:` + ActionUnstack
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
	Commands[CmdContainerPropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/container/` + PlHoldTomID + `/property/`,
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
	Commands[CmdContainerRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/container/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerResolve] = CmdDef{
		Method:      MethodGET,
		Path:        `/container/` + PlHoldTomID + `/resolve/` + PlHoldResolv,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{PlHoldTomID, PlHoldResolv},
	}
	Commands[CmdContainerShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/container/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerStack] = CmdDef{
		Method:      MethodPUT,
		Path:        `/container/` + PlHoldTomID + `/parent`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdContainerUnstack] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/container/` + PlHoldTomID + `/parent`,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
}

// Container ...
type Container struct {
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Type      string                    `json:"type"`
	Parent    string                    `json:"parent"`
	Link      []string                  `json:"link"`
	Property  map[string]PropertyDetail `json:"property"`
	CreatedAt string                    `json:"createdAt,omitempty"`
	CreatedBy string                    `json:"createdBy,omitempty"`
	Resources []string                  `json:"resources"`
	ID        string                    `json:"-"`
	TomID     string                    `json:"-"`
}

func NewContainerRequest() Request {
	return Request{
		Container: NewContainer(),
	}
}

func NewContainer() *Container {
	return &Container{
		Link:      []string{},
		Property:  map[string]PropertyDetail{},
		Resources: []string{},
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

func (c *Container) ExportName() string {
	return c.Name
}

func (c *Container) ExportNamespace() string {
	return c.Namespace
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
