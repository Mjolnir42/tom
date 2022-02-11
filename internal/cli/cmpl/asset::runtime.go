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

func RuntimeAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeAdd)
	GenericMultiWithProperty(c, uniq, multi)
}

func RuntimeRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeRemove)
	GenericMulti(c, uniq, multi)
}

func RuntimeList(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeList)
	DirectMulti(c, uniq, multi)
}

func RuntimeShow(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeShow)
	GenericMulti(c, uniq, multi)
}

func RuntimePropSet(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimePropSet)
	GenericPropertyChain(c, uniq, multi)
}

func RuntimePropUpdate(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimePropUpdate)
	GenericPropertyChain(c, uniq, multi)
}

func RuntimePropRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimePropRemove)
	GenericMulti(c, uniq, multi)
}

func RuntimeLink(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeLink)
	GenericMulti(c, uniq, multi)
}

func RuntimeStack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeStack)
	GenericMulti(c, uniq, multi)
}

func RuntimeUnstack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeUnstack)
	GenericMulti(c, uniq, multi)
}

func RuntimeResolve(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeResolve)
	GenericMulti(c, uniq, multi)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
