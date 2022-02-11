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
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerRemove)
	GenericMulti(c, uniq, multi)
}

func ServerList(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerList)
	DirectMulti(c, uniq, multi)
}

func ServerShow(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerShow)
	GenericMulti(c, uniq, multi)
}

func ServerPropSet(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerPropSet)
	GenericPropertyChain(c, uniq, multi)
}

func ServerPropUpdate(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerPropUpdate)
	GenericPropertyChain(c, uniq, multi)
}

func ServerPropRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerPropRemove)
	GenericMulti(c, uniq, multi)
}

func ServerLink(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerLink)
	GenericMulti(c, uniq, multi)
}

func ServerStack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerStack)
	GenericMulti(c, uniq, multi)
}

func ServerUnstack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdServerUnstack)
	GenericMulti(c, uniq, multi)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
