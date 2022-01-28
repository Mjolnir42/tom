/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/cli/model/asset"

import (
	"github.com/mjolnir42/tom/internal/cli/cmpl"
	"github.com/mjolnir42/tom/internal/cli/help"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func registerAssetOrchestration(app cli.App, run ActionFunc) *cli.App {
	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `orchestration`,
				Usage:       `Commands for maintaining orchestration environments`,
				Description: help.Text(proto.CmdOrchestration),
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new orchestration environment`,
						Description:  help.Text(proto.CmdOrchestrationAdd),
						Action:       run(cmdAssetOrchestrationAdd),
						BashComplete: cmpl.OrchestrationAdd,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
