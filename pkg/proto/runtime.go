/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Runtime defines a runtime within the asset model
type Runtime struct {
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

func NewRuntimeRequest() Request {
	return Request{
		Runtime: &Runtime{},
	}
}

// RuntimeHeader defines ....
type RuntimeHeader struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
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
	var ntt string
	switch {
	case r.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(r.TomID):
		r.Name, r.Namespace, ntt = parseTomIDFormatDNS(r.TomID)
		if ntt != EntityRuntime {
			return ErrInvalidTomID
		}
		return nil
	case isTomIDFormatURI(r.TomID):
		r.Name, r.Namespace, ntt = parseTomIDFormatURI(r.TomID)
		if ntt != EntityRuntime {
			return ErrInvalidTomID
		}
		return nil
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
