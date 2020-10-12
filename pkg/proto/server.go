/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const EntityServer = `server`

// Server defines a server within the asset model
type Server struct {
	ID        string     `json:"-"`
	TomID     string     `json:"-"`
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Parent    string     `json:"parent"`
	Link      []string   `json:"link"`
	Property  []Property `json:"property"`
}

func (s *Server) String() string {
	return s.DNS()
}

func (s *Server) DNS() string {
	return s.Name + `.` + s.Namespace + `.` + EntityServer + `.tom`
}

func (s *Server) TomID() string {
	return `tom://` + s.Namespace + `/` + EntityServer + `/name=` + s.Name
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
