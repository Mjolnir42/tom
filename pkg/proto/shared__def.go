/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"fmt"
	"strings"
)

type CmdDef struct {
	Method      string
	Path        string
	Body        bool
	ResultTmpl  string
	Placeholder []string
}

var Commands = map[string]CmdDef{
	CmdNamespaceAdd: {
		Method:      MethodPOST,
		Path:        `/namespace/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	},
	CmdNamespaceList: {
		Method:      MethodGET,
		Path:        `/namespace/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	},
	CmdNamespaceShow: {
		Method:      MethodGET,
		Path:        `/namespace/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	},
	CmdNamespaceAttrAdd: {
		Method:      MethodPOST,
		Path:        `/namespace/` + PlHoldTomID + `/attribute/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
}

func AssertCommandIsDefined(c string) {
	if _, ok := Commands[c]; !ok {
		panic(c)
	}
}

// OnlyUnreserved checks that s only contains charcters from proto.CharUnreserved
func OnlyUnreserved(s string) error {
	for b := range []byte(s) {
		if !strings.Contains(CharUnreserved, string(b)) {
			return fmt.Errorf(
				"String <%s> contains illegal character <%s>",
				s,
				string(b),
			)
		}
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
