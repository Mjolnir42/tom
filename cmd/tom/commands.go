/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"github.com/mjolnir42/tom/internal/cli/model"
	"github.com/mjolnir42/tom/internal/cli/model/asset"
	"github.com/mjolnir42/tom/internal/cli/model/meta"
	"github.com/urfave/cli/v2"
)

var Registry map[string]func(*cli.Context) error

func init() {
	Registry = make(map[string]func(*cli.Context) error)
}

func registerCommands(app cli.App) *cli.App {

	app = *meta.Register(app, runtime, Registry)
	app = *asset.Register(app, runtime, Registry)
	app = *model.Register(app, runtime, Registry)

	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
