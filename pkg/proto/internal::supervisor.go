/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdSupervisor        = ModelInternal + `::` + EntitySupervisor + `:`
	CmdSupervisorAuthEPK = ModelInternal + `::` + EntitySupervisor + `:` + ActionAuthenticateEPK
)

func init() {
	Commands[CmdSupervisorAuthEPK] = CmdDef{
		Method:      MethodHEAD,
		Path:        `/supervisor`,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
