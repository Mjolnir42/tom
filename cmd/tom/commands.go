/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/mjolnir42/tom/internal/cli/help"
	"github.com/mjolnir42/tom/internal/cli/model"
	"github.com/mjolnir42/tom/internal/cli/model/asset"
	"github.com/mjolnir42/tom/internal/cli/model/bulk"
	"github.com/mjolnir42/tom/internal/cli/model/meta"
	"github.com/urfave/cli/v2"
)

var Registry map[string]func(*cli.Context) error

func init() {
	Registry = make(map[string]func(*cli.Context) error)
}

func registerCommands(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:   `output-autocomplete`,
				Hidden: true,
				Action: cmdOutputAutoComplete,
			},
		}...,
	)
	app.Flags = append(app.Flags,
		[]cli.Flag{
			&cli.BoolFlag{
				Name:    `verbose`,
				Aliases: []string{`v`},
				Usage:   `request verbose service replies`,
				Value:   false,
			},
		}...,
	)

	app = *meta.Register(app, runtime, Registry)
	app = *asset.Register(app, runtime, Registry)
	app = *bulk.Register(app, runtime, Registry)
	app = *model.Register(app, runtime, Registry)

	return &app
}

func cmdOutputAutoComplete(c *cli.Context) error {
	fmt.Println(help.Text(`zsh_autocomplete`))
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
