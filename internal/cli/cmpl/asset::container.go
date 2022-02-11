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

func ContainerAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdRuntimeAdd)
	GenericMultiWithProperty(c, uniq, multi)
}

func ContainerRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerRemove)
	GenericMulti(c, uniq, multi)
}

func ContainerList(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerList)
	DirectMulti(c, uniq, multi)
}

func ContainerShow(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerShow)
	GenericMulti(c, uniq, multi)
}

func ContainerPropSet(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerPropSet)
	GenericPropertyChain(c, uniq, multi)
}

func ContainerPropUpdate(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerPropUpdate)
	GenericPropertyChain(c, uniq, multi)
}

func ContainerPropRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerPropRemove)
	GenericPropertyChain(c, uniq, multi)
}

func ContainerLink(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerLink)
	GenericMulti(c, uniq, multi)
}

func ContainerStack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerStack)
	GenericMulti(c, uniq, multi)
}

func ContainerUnstack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerUnstack)
	GenericMulti(c, uniq, multi)
}

func ContainerResolve(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdContainerResolve)
	GenericMulti(c, uniq, multi)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
