/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// CmdDef holds required information for a specific command, allowing to
// make a request for this command
type CmdDef struct {
	Method      string
	Path        string
	Body        bool
	ResultTmpl  string
	Placeholder []string
}

// Commands is the exported map of all protocol commands defined in proto
var Commands = map[string]CmdDef{}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
