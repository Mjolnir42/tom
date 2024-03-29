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

func registerAssetRuntime(app cli.App, run ActionFunc) *cli.App {
	handlerMap[proto.EntityRuntime+`:`+proto.ActionList] = run(cmdAssetRuntimeList)
	handlerMap[proto.EntityRuntime+`:`+proto.ActionShow] = run(cmdAssetRuntimeShow)

	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `runtime`,
				Usage:       `Commands for maintaining runtime environments`,
				Description: help.Text(proto.CmdRuntime),
				Flags:       []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new runtime environment`,
						Description:  help.Text(proto.CmdRuntimeAdd),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeAdd),
						BashComplete: cmpl.RuntimeAdd,
					},
					{
						Name:         `list`,
						Usage:        `List all existing runtime environments`,
						Description:  help.Text(proto.CmdRuntimeList),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeList),
						BashComplete: cmpl.RuntimeList,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a runtime environment`,
						Description:  help.Text(proto.CmdRuntimeShow),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeShow),
						BashComplete: cmpl.RuntimeShow,
					},
					{
						Name:        `property`,
						Usage:       `Commands for runtime environment properties`,
						Description: help.Text(proto.CmdRuntime),
						Flags:       []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Subcommands: []*cli.Command{
							{
								Name:         `set`,
								Usage:        `Set all properties of a runtime environment`,
								Description:  help.Text(proto.CmdRuntimePropSet),
								Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
								Action:       run(cmdAssetRuntimePropSet),
								BashComplete: cmpl.RuntimePropSet,
							},
							{
								Name:         `update`,
								Usage:        `Update properties of a runtime environment`,
								Description:  help.Text(proto.CmdRuntimePropUpdate),
								Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
								Action:       run(cmdAssetRuntimePropUpdate),
								BashComplete: cmpl.RuntimePropUpdate,
							},
							{
								Name:         `remove`,
								Usage:        `Remove properties of a runtime environment`,
								Description:  help.Text(proto.CmdRuntimePropRemove),
								Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
								Action:       run(cmdAssetRuntimePropRemove),
								BashComplete: cmpl.RuntimePropRemove,
							},
						},
					},
					{
						Name:         `remove`,
						Usage:        `Remove a runtime environment`,
						Description:  help.Text(proto.CmdRuntimeRemove),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeRemove),
						BashComplete: cmpl.RuntimeRemove,
					},
					{
						Name:         `link`,
						Usage:        `Link two runtime environments as referring to the same entity`,
						Description:  help.Text(proto.CmdRuntimeLink),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeLink),
						BashComplete: cmpl.RuntimeLink,
					},
					{
						Name:         `stack`,
						Usage:        `Specify the base this runtime runs on`,
						Description:  help.Text(proto.CmdRuntimeStack),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeStack),
						BashComplete: cmpl.RuntimeStack,
					},
					{
						Name:         `unstack`,
						Usage:        `Remove the base this runtime runs on`,
						Description:  help.Text(proto.CmdRuntimeUnstack),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeUnstack),
						BashComplete: cmpl.RuntimeUnstack,
					},
					{
						Name:         `resolve`,
						Usage:        `Resolve this runtime down to the server(s) it runs on`,
						Description:  help.Text(proto.CmdRuntimeResolve),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdAssetRuntimeResolve),
						BashComplete: cmpl.RuntimeResolve,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
