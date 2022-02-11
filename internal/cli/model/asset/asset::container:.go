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

func registerAssetContainer(app cli.App, run ActionFunc) *cli.App {
	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `container`,
				Usage:       `Commands for maintaining containers`,
				Description: help.Text(proto.CmdContainer),
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new container`,
						Description:  help.Text(proto.CmdContainerAdd),
						Action:       run(cmdAssetContainerAdd),
						BashComplete: cmpl.ContainerAdd,
					},
					{
						Name:         `list`,
						Usage:        `List all existing containers`,
						Description:  help.Text(proto.CmdContainerList),
						Action:       run(cmdAssetContainerList),
						BashComplete: cmpl.ContainerList,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a container`,
						Description:  help.Text(proto.CmdContainerShow),
						Action:       run(cmdAssetContainerShow),
						BashComplete: cmpl.ContainerShow,
					},
					{
						Name:        `property`,
						Usage:       `Commands for container properties`,
						Description: help.Text(proto.CmdContainer),
						Subcommands: []*cli.Command{
							{
								Name:         `set`,
								Usage:        `Set all properties of a container`,
								Description:  help.Text(proto.CmdContainerPropSet),
								Action:       run(cmdAssetContainerPropSet),
								BashComplete: cmpl.ContainerPropSet,
							},
							{
								Name:         `update`,
								Usage:        `Update properties of a container`,
								Description:  help.Text(proto.CmdContainerPropUpdate),
								Action:       run(cmdAssetContainerPropUpdate),
								BashComplete: cmpl.ContainerPropUpdate,
							},
							{
								Name:         `remove`,
								Usage:        `Remove properties of a container`,
								Description:  help.Text(proto.CmdContainerPropRemove),
								Action:       run(cmdAssetContainerPropRemove),
								BashComplete: cmpl.ContainerPropRemove,
							},
						},
					},
					{
						Name:         `remove`,
						Usage:        `Remove a container`,
						Description:  help.Text(proto.CmdContainerRemove),
						Action:       run(cmdAssetContainerRemove),
						BashComplete: cmpl.ContainerRemove,
					},
					{
						Name:         `link`,
						Usage:        `Link two containers as referring to the same entity`,
						Description:  help.Text(proto.CmdContainerLink),
						Action:       run(cmdAssetContainerLink),
						BashComplete: cmpl.ContainerLink,
					},
					{
						Name:         `stack`,
						Usage:        `Specify the runtime environment providing for a container`,
						Description:  help.Text(proto.CmdContainerStack),
						Action:       run(cmdAssetContainerStack),
						BashComplete: cmpl.ContainerStack,
					},
					{
						Name:         `unstack`,
						Usage:        `Remove the provider runtime environment of the container`,
						Description:  help.Text(proto.CmdContainerUnstack),
						Action:       run(cmdAssetContainerUnstack),
						BashComplete: cmpl.ContainerUnstack,
					},
					{
						Name:         `resolve`,
						Usage:        `Resolve this container down to the server(s) it runs on`,
						Description:  help.Text(proto.CmdContainerResolve),
						Action:       run(cmdAssetContainerResolve),
						BashComplete: cmpl.ContainerResolve,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
