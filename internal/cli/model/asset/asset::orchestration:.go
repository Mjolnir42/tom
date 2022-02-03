/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
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
					{
						Name:         `list`,
						Usage:        `List all existing orchestration environments`,
						Description:  help.Text(proto.CmdOrchestrationList),
						Action:       run(cmdAssetOrchestrationList),
						BashComplete: cmpl.OrchestrationList,
					},
					{
						Name:         `show`,
						Usage:        `Show full details about an orchestration environment`,
						Description:  help.Text(proto.CmdOrchestrationShow),
						Action:       run(cmdAssetOrchestrationShow),
						BashComplete: cmpl.OrchestrationShow,
					},
					{
						Name:        `property`,
						Usage:       `Commands for orchestration environment properties`,
						Description: help.Text(proto.CmdOrchestration),
						Subcommands: []*cli.Command{
							{
								Name:         `set`,
								Usage:        `Set all properties of an orchestration environment`,
								Description:  help.Text(proto.CmdOrchestrationPropSet),
								Action:       run(cmdAssetOrchestrationPropSet),
								BashComplete: cmpl.OrchestrationPropSet,
							},
							{
								Name:         `update`,
								Usage:        `Update properties of an orchestration environment`,
								Description:  help.Text(proto.CmdOrchestrationPropUpdate),
								Action:       run(cmdAssetOrchestrationPropUpdate),
								BashComplete: cmpl.OrchestrationPropUpdate,
							},
						},
					},
					{
						Name:         `link`,
						Usage:        `Link two orchestration environments as referring to the same entity`,
						Description:  help.Text(proto.CmdOrchestrationLink),
						Action:       run(cmdAssetOrchestrationLink),
						BashComplete: cmpl.OrchestrationLink,
					},
					{
						Name:         `stack`,
						Usage:        `Specify base runtimes an orchestration sits on top of`,
						Description:  help.Text(proto.CmdOrchestrationStack),
						Action:       run(cmdAssetOrchestrationStack),
						BashComplete: cmpl.OrchestrationStack,
					},
					{
						Name:         `unstack`,
						Usage:        `Remove base runtimes from an orchestration`,
						Description:  help.Text(proto.CmdOrchestrationUnstack),
						Action:       run(cmdAssetOrchestrationUnstack),
						BashComplete: cmpl.OrchestrationUnstack,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
