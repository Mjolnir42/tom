/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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

// OnlyUnreserved checks that s only contains characters from proto.CharUnreserved
func OnlyUnreserved(s string) error {
	for _, b := range []byte(s) {
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

// ValidNamespace checks that s only contains characters and structure suitable
// for a namespace name
func ValidNamespace(s string) error {
	for _, b := range []byte(s) {
		if !strings.Contains(CharNamespace, string(b)) {
			return fmt.Errorf(
				"String <%s> contains illegal character <%s>",
				s,
				string(b),
			)
		}
	}
	switch strings.Count(s, `~`) {
	case 0:
	case 1:
		switch {
		case strings.HasPrefix(s, `tool~`):
		case strings.HasPrefix(s, `team~`):
		default:
			return fmt.Errorf(
				"Unknown namespace prefix: %s",
				strings.Split(s, `~`)[1],
			)
		}
	default:
		return fmt.Errorf(
			"Character <~> matches %d times but is only allowed once as prefix separator",
			strings.Count(s, `~`),
		)
	}
	return nil
}

// CheckPropertyConstraints checks that if the property attribute end
// with one of the suffixes _list, _json or _xml, then the value must
// decode appropriately.
func CheckPropertyConstraints(prop *PropertyDetail) error {
	if strings.HasSuffix(prop.Attribute, `_list`) {
		jl := &[]string{}
		if err := json.NewDecoder(bytes.NewBufferString(prop.Value)).Decode(jl); err != nil {
			return fmt.Errorf("Property %s is not a valid list: %s",
				prop.Attribute,
				err.Error(),
			)
		}
		return nil
	}
	if strings.HasSuffix(prop.Attribute, `_json`) {
		var j interface{}
		if err := json.NewDecoder(bytes.NewBufferString(prop.Value)).Decode(&j); err != nil {
			return fmt.Errorf("Property %s is not valid JSON: %s",
				prop.Attribute,
				err.Error(),
			)
		}
		return nil
	}
	if strings.HasSuffix(prop.Attribute, `_xml`) {
		var x interface{}
		if err := xml.Unmarshal([]byte(prop.Value), &x); err != nil {
			return fmt.Errorf("Property %s is not supported XML: %s",
				prop.Attribute,
				err.Error(),
			)
		}
		return nil
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
