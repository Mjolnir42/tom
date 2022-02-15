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

func registerAssetServer(app cli.App, run ActionFunc) *cli.App {
	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `server`,
				Usage:       `Commands for maintaining servers`,
				Description: help.Text(proto.CmdServer),
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new server`,
						Description:  help.Text(proto.CmdServerAdd),
						Action:       run(cmdAssetServerAdd),
						BashComplete: cmpl.ServerAdd,
					},
					{
						Name:         `list`,
						Usage:        `List all existing servers`,
						Description:  help.Text(proto.CmdServerList),
						Action:       run(cmdAssetServerList),
						BashComplete: cmpl.ServerList,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a servers`,
						Description:  help.Text(proto.CmdServerShow),
						Action:       run(cmdAssetServerShow),
						BashComplete: cmpl.ServerShow,
					},
					{
						Name:        `property`,
						Usage:       `Commands for server properties`,
						Description: help.Text(proto.CmdServer),
						Subcommands: []*cli.Command{
							{
								Name:         `set`,
								Usage:        `Set all properties of a server`,
								Description:  help.Text(proto.CmdServerPropSet),
								Action:       run(cmdAssetServerPropSet),
								BashComplete: cmpl.ServerPropSet,
							},
							{
								Name:         `update`,
								Usage:        `Update properties of a server`,
								Description:  help.Text(proto.CmdServerPropUpdate),
								Action:       run(cmdAssetServerPropUpdate),
								BashComplete: cmpl.ServerPropUpdate,
							},
							{
								Name:         `remove`,
								Usage:        `Remove properties of a server`,
								Description:  help.Text(proto.CmdServerPropRemove),
								Action:       run(cmdAssetServerPropRemove),
								BashComplete: cmpl.ServerPropRemove,
							},
						},
					},
					{
						Name:         `remove`,
						Usage:        `Remove a server`,
						Description:  help.Text(proto.CmdServerRemove),
						Action:       run(cmdAssetServerRemove),
						BashComplete: cmpl.ServerRemove,
					},
					{
						Name:         `link`,
						Usage:        `Link two servers as referring to the same entity`,
						Description:  help.Text(proto.CmdServerLink),
						Action:       run(cmdAssetServerLink),
						BashComplete: cmpl.ServerLink,
					},
					{
						Name:         `stack`,
						Usage:        `Specify the runtime environment providing a server`,
						Description:  help.Text(proto.CmdServerStack),
						Action:       run(cmdAssetServerStack),
						BashComplete: cmpl.ServerStack,
					},
					{
						Name:         `unstack`,
						Usage:        `Remove the provider runtime environment of a server`,
						Description:  help.Text(proto.CmdServerUnstack),
						Action:       run(cmdAssetServerUnstack),
						BashComplete: cmpl.ServerUnstack,
					},
					{
						Name:         `resolve`,
						Usage:        `Resolve this server down to the server(s) it runs on`,
						Description:  help.Text(proto.CmdServerResolve),
						Action:       run(cmdAssetServerResolve),
						BashComplete: cmpl.ServerResolve,
					},
					{
						Name:         `remove`,
						Usage:        `Remove a server`,
						Description:  help.Text(proto.CmdServerRemove),
						Action:       run(cmdAssetServerRemove),
						BashComplete: cmpl.ServerRemove,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
