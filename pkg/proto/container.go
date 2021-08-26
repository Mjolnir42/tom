/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Container ...
type Container struct {
	ID           string                    `json:"-"`
	TomID        string                    `json:"-"`
	Namespace    string                    `json:"namespace"`
	Name         string                    `json:"name"`
	Type         string                    `json:"type"`
	Parent       string                    `json:"parent"`
	Link         []string                  `json:"link"`
	Property     map[string]PropertyDetail `json:"property"`
	StdProperty  []PropertyDetail          `json:"-"`
	UniqProperty []PropertyDetail          `json:"-"`
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
	switch {
	case c.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(c.TomID):
		c.Name, c.Namespace, _ = parseTomIDFormatDNS(c.TomID)
		return nil
	case isTomIDFormatURI(c.TomID):
		c.Name, c.Namespace, _ = parseTomIDFormatURI(c.TomID)
		return nil
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
