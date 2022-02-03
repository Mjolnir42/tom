/*-
 * Copyright (c) 2022, Jörg Pernfuß <joerg.pernfuss@ionos.com>
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

func OrchestrationAdd(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationAdd)
	GenericMultiWithProperty(c, uniq, multi)
}

func OrchestrationRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationRemove)
	GenericMulti(c, uniq, multi)
}

func OrchestrationList(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationList)
	DirectMulti(c, uniq, multi)
}

func OrchestrationShow(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationShow)
	GenericMulti(c, uniq, multi)
}

func OrchestrationPropSet(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationPropSet)
	GenericPropertyChain(c, uniq, multi)
}

func OrchestrationPropUpdate(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationPropUpdate)
	GenericPropertyChain(c, uniq, multi)
}

func OrchestrationPropRemove(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationPropRemove)
	GenericMulti(c, uniq, multi)
}

func OrchestrationLink(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationLink)
	GenericMulti(c, uniq, multi)
}

func OrchestrationStack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationStack)
	GenericMulti(c, uniq, multi)
}

func OrchestrationUnstack(c *cli.Context) {
	multi, uniq, _ := adm.ArgumentsForCommand(proto.CmdOrchestrationUnstack)
	GenericMulti(c, uniq, multi)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
