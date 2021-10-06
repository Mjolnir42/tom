/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/cli/model/asset"

import (
	"github.com/urfave/cli/v2"
)

type ActionFunc func(cli.ActionFunc) cli.ActionFunc

func Register(app cli.App, run ActionFunc) *cli.App {
	app = *registerAssetServer(app, run)
	app = *registerAssetRuntime(app, run)
	app = *registerAssetContainer(app, run)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
