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
	GenericMulti(c, uniq, multi)
}

func RuntimeRemove(c *cli.Context) {
	None(c)
}

func RuntimeList(c *cli.Context) {
	None(c)
}

func RuntimeShow(c *cli.Context) {
	None(c)
}

func RuntimePropSet(c *cli.Context) {
	None(c)
}

func RuntimePropUpdate(c *cli.Context) {
	None(c)
}

func RuntimePropRemove(c *cli.Context) {
	None(c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
