/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

type CmdDef struct {
	Method string
	Path   string
	Body   bool
}

var Commands = map[string]CmdDef{
	CmdNamespaceAdd: {
		Method: MethodPOST,
		Path:   `/namespace/`,
		Body:   true,
	},
	CmdNamespaceList: {
		Method: MethodGET,
		Path:   `/namespace/`,
		Body:   false,
	},
	CmdNamespaceShow: {
		Method: MethodGET,
		Path:   `/namespace/:tomID`,
		Body:   false,
	},
}

func AssertCommandIsDefined(c string) {
	if _, ok := Commands[c]; !ok {
		panic(c)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
