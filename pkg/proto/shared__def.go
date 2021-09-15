/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

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
	CmdNamespaceAttrRemove: {
		Method:      MethodDELETE,
		Path:        `/namespace/` + PlHoldTomID + `/attribute/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdNamespacePropSet: {
		Method:      MethodPUT,
		Path:        `/namespace/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdNamespacePropUpdate: {
		Method:      MethodPATCH,
		Path:        `/namespace/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdNamespacePropRemove: {
		Method:      MethodDELETE,
		Path:        `/namespace/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdNamespaceRemove: {
		Method:      MethodDELETE,
		Path:        `/namespace/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdServerList: {
		Method:      MethodGET,
		Path:        `/server/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	},
	CmdServerShow: {
		Method:      MethodGET,
		Path:        `/server/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	},
	CmdServerAdd: {
		Method:      MethodPOST,
		Path:        `/server/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	},
	CmdServerRemove: {
		Method:      MethodDELETE,
		Path:        `/server/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdServerPropSet: {
		Method:      MethodPUT,
		Path:        `/server/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdServerPropUpdate: {
		Method:      MethodPATCH,
		Path:        `/server/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdServerPropRemove: {
		Method:      MethodDELETE,
		Path:        `/server/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdRuntimeList: {
		Method:      MethodGET,
		Path:        `/runtime/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	},
	CmdRuntimeShow: {
		Method:      MethodGET,
		Path:        `/runtime/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	},
	CmdRuntimeAdd: {
		Method:      MethodPOST,
		Path:        `/runtime/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	},
	CmdRuntimeRemove: {
		Method:      MethodDELETE,
		Path:        `/runtime/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdRuntimePropSet: {
		Method:      MethodPUT,
		Path:        `/runtime/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdRuntimePropUpdate: {
		Method:      MethodPATCH,
		Path:        `/runtime/` + PlHoldTomID + `/property/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	},
	CmdRuntimePropRemove: {
		Method:      MethodDELETE,
		Path:        `/runtime/` + PlHoldTomID + `/property/`,
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
