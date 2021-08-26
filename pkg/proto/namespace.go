/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Namespace defines ...
type Namespace struct {
	Name         string                    `json:"name"`
	Type         string                    `json:"type"`
	LookupKey    string                    `json:"lookup-attribute-key"`
	LookupURI    string                    `json:"lookup-uri"`
	Constraint   []string                  `json:"entity-contraint-list"`
	Attributes   []AttributeDefinition     `json:"attributes"`
	Property     map[string]PropertyDetail `json:"property,omitempty"`
	CreatedAt    string                    `json:"createdAt"`
	CreatedBy    string                    `json:"createdBy"`
	ID           string                    `json:"-"`
	TomID        string                    `json:"-"`
	StdProperty  []PropertyDetail          `json:"-"`
	UniqProperty []PropertyDetail          `json:"-"`
}

func NewNamespaceRequest() Request {
	return Request{
		Namespace: &Namespace{},
	}
}

// NamespaceHeader defines ...
type NamespaceHeader struct {
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
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
	switch {
	case n.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(n.TomID):
		n.Name, _, _ = parseTomIDFormatDNS(n.TomID)
		return nil
	case isTomIDFormatURI(n.TomID):
		n.Name, _, _ = parseTomIDFormatURI(n.TomID)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
