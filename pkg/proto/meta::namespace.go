/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdNamespace           = ModelMeta + `::` + EntityNamespace + `:`
	CmdNamespaceAdd        = ModelMeta + `::` + EntityNamespace + `:` + ActionAdd
	CmdNamespaceAttrAdd    = ModelMeta + `::` + EntityNamespace + `:` + ActionAttrAdd
	CmdNamespaceAttrRemove = ModelMeta + `::` + EntityNamespace + `:` + ActionAttrRemove
	CmdNamespaceList       = ModelMeta + `::` + EntityNamespace + `:` + ActionList
	CmdNamespacePropRemove = ModelMeta + `::` + EntityNamespace + `:` + ActionPropRemove
	CmdNamespacePropSet    = ModelMeta + `::` + EntityNamespace + `:` + ActionPropSet
	CmdNamespacePropUpdate = ModelMeta + `::` + EntityNamespace + `:` + ActionPropUpdate
	CmdNamespaceRemove     = ModelMeta + `::` + EntityNamespace + `:` + ActionRemove
	CmdNamespaceShow       = ModelMeta + `::` + EntityNamespace + `:` + ActionShow
)

func init() {
	Commands[CmdNamespaceAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/namespace/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdNamespaceList] = CmdDef{
		Method:      MethodGET,
		Path:        `/namespace/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdNamespaceShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/namespace/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdNamespaceAttrAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/namespace/` + PlHoldTomID + `/attribute/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdNamespaceAttrRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/namespace/` + PlHoldTomID + `/attribute/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdNamespacePropSet] = CmdDef{
		Method:      MethodPUT,
		Path:        `/namespace/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdNamespacePropUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/namespace/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdNamespacePropRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/namespace/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdNamespaceRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/namespace/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
}

// Namespace defines ...
type Namespace struct {
	Name       string                    `json:"name"`
	Type       string                    `json:"type"`
	LookupKey  string                    `json:"lookup-attribute-key"`
	LookupURI  string                    `json:"lookup-uri"`
	Constraint []string                  `json:"entity-contraint-list"`
	Attributes []AttributeDefinition     `json:"attributes"`
	Property   map[string]PropertyDetail `json:"property,omitempty"`
	CreatedAt  string                    `json:"createdAt,omitempty"`
	CreatedBy  string                    `json:"createdBy,omitempty"`
	ID         string                    `json:"-"`
	TomID      string                    `json:"-"`
}

func NewNamespaceRequest() Request {
	return Request{
		Namespace: NewNamespace(),
	}
}

func NewNamespace() *Namespace {
	return &Namespace{
		Constraint: []string{},
		Attributes: []AttributeDefinition{},
		Property:   map[string]PropertyDetail{},
	}
}

// NamespaceHeader defines ...
type NamespaceHeader struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

func (n *Namespace) SetTomID() Entity {
	n.TomID = n.FormatDNS()
	return n
}

func (n *Namespace) String() string {
	return n.FormatDNS()
}

func (n *Namespace) FormatDNS() string {
	return n.Name + `.` + EntityNamespace + `.tom`
}

func (n *Namespace) FormatTomID() string {
	return `tom:///` + EntityNamespace + `/name=` + n.Name
}

func (n *Namespace) ParseTomID() error {
	var ntt string
	switch {
	case n.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(n.TomID):
		n.Name, _, ntt = parseTomIDFormatDNS(n.TomID)
		if ntt != EntityNamespace {
			return ErrInvalidTomID
		}
		return nil
	case isTomIDFormatURI(n.TomID):
		n.Name, _, ntt = parseTomIDFormatURI(n.TomID)
		if ntt != EntityNamespace {
			return ErrInvalidTomID
		}
		return nil
	default:
		return ErrInvalidTomID
	}
}

func (n *Namespace) PropertyIterator() <-chan PropertyDetail {
	ret := make(chan PropertyDetail)
	go func() {
		for key := range n.Property {
			ret <- n.Property[key]
		}
		close(ret)
	}()
	return ret
}

func (n *Namespace) ExportName() string {
	return n.Name
}

func (n *Namespace) ExportNamespace() string {
	return n.Name
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
