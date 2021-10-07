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

func ServerAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerAdd)
	GenericMultiWithProperty(c, uniq, multi)
}

func ServerRemove(c *cli.Context) {
	None(c)
}

func ServerList(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerList)
	DirectMulti(c, uniq, multi)
}

func ServerShow(c *cli.Context) {
	None(c)
}

func ServerPropSet(c *cli.Context) {
	None(c)
}

func ServerPropUpdate(c *cli.Context) {
	None(c)
}

func ServerPropRemove(c *cli.Context) {
	None(c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
