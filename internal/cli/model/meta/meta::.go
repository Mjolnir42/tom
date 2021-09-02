/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/cli/model/meta"

import (
	"github.com/urfave/cli/v2"
)

type Runtime func(cli.ActionFunc) cli.ActionFunc

func Register(app cli.App, run Runtime) *cli.App {
	return registerMetaNamespace(app, run)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
