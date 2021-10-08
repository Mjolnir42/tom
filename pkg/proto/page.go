/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto // import "github.com/mjolnir42/tom/pkg/proto/"

func init() {
	Commands[CmdPageStaticCSS] = CmdDef{
		Method:      MethodGET,
		Path:        `/static/css/` + PlHoldAsset,
		Body:        false,
		ResultTmpl:  TemplateBindata,
		Placeholder: []string{PlHoldAsset},
	}
	Commands[CmdPageStaticFont] = CmdDef{
		Method:      MethodGET,
		Path:        `/static/fonts/` + PlHoldAsset,
		Body:        false,
		ResultTmpl:  TemplateBindata,
		Placeholder: []string{PlHoldAsset},
	}
	Commands[CmdPageStaticImage] = CmdDef{
		Method:      MethodGET,
		Path:        `/static/img/` + PlHoldAsset,
		Body:        false,
		ResultTmpl:  TemplateBindata,
		Placeholder: []string{PlHoldAsset},
	}
	Commands[CmdPageStaticJS] = CmdDef{
		Method:      MethodGET,
		Path:        `/static/js/` + PlHoldAsset,
		Body:        false,
		ResultTmpl:  TemplateBindata,
		Placeholder: []string{PlHoldAsset},
	}
	Commands[CmdPageApplication] = CmdDef{
		Method:      MethodGET,
		Path:        `/`,
		Body:        false,
		ResultTmpl:  TemplateNone,
		Placeholder: []string{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
