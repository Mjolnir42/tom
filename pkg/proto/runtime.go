/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const EntityRuntime = `runtime`

// Runtime defines a runtime within the asset model
type Runtime struct {
	ID           string            `json:"-"`
	TomID        string            `json:"-"`
	Namespace    string            `json:"namespace"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Parent       string            `json:"parent"`
	Link         []string          `json:"link"`
	PropertyMap  map[string]string `json:"property"`
	StdProperty  []Property        `json:"-"`
	UniqProperty []Property        `json:"-"`
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
	switch {
	case r.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(r.TomID):
		r.Name, r.Namespace, _ = parseTomIDFormatDNS(r.TomID)
		return nil
	case isTomIDFormatURI(r.TomID):
		r.Name, r.Namespace, _ = parseTomIDFormatURI(r.TomID)
		return nil
	default:
		return ErrInvalidTomID
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
