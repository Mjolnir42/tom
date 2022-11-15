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

func FlowAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowAdd)
	GenericMultiWithProperty(c, uniq, multi)
}

func FlowEnsure(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowEnsure)
	GenericMultiWithProperty(c, uniq, multi)
}

func FlowRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowRemove)
	GenericMulti(c, uniq, multi)
}

func FlowList(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowList)
	DirectMulti(c, uniq, multi)
}

func FlowShow(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowShow)
	GenericMulti(c, uniq, multi)
}

func FlowPropSet(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowPropSet)
	GenericPropertyChain(c, uniq, multi)
}

func FlowPropUpdate(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowPropUpdate)
	GenericPropertyChain(c, uniq, multi)
}

func FlowPropRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdFlowPropRemove)
	GenericPropertyChain(c, uniq, multi)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
