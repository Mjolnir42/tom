/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/cli/model/meta"

import (
	"github.com/mjolnir42/tom/internal/cli/cmpl"
	"github.com/mjolnir42/tom/internal/cli/help"
	"github.com/urfave/cli/v2"
)

func registerMetaNamespace(app cli.App, run Runtime) *cli.App {
	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `namespace`,
				Usage:       `Commands for maintaining namespaces`,
				Description: help.Text(`meta::namespace:`),
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new namespace`,
						Description:  help.Text(`meta::namespace:add`),
						Action:       run(cmdMetaNamespaceAdd),
						BashComplete: cmpl.NamespaceAdd,
					},
					{
						Name:         `list`,
						Usage:        `List all existing namespaces`,
						Description:  help.Text(`meta::namespace:list`),
						Action:       run(cmdMetaNamespaceList),
						BashComplete: cmpl.NamespaceList,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
