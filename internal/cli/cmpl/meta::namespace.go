/*-
 * Copyright (c) 2021, Jörg Pernfuß <joerg.pernfuss@ionos.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cmpl

import (
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func NamespaceAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdNamespaceAdd)
	GenericMulti(c, uniq, multi)
}

func NamespaceList(c *cli.Context) {
	None(c)
}

func NamespaceShow(c *cli.Context) {
	None(c)
}

func NamespaceAttrAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdNamespaceAttrAdd)
	GenericMulti(c, uniq, multi)
}

func NamespaceAttrRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdNamespaceAttrRemove)
	GenericMulti(c, uniq, multi)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
