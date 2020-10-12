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
	ID        string     `json:"-"`
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Parent    string     `json:"parent"`
	Link      []string   `json:"link"`
	Property  []Property `json:"property"`
}

func (r *Runtime) String() string {
	return r.DNS()
}

func (r *Runtime) DNS() string {
	return r.Name + `.` + r.Namespace + `.` + EntityRuntime + `.tom`
}

func (r *Runtime) TomID() string {
	return `tom://` + r.Namespace + `/` + EntityRuntime + `/name=` + r.Name
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
