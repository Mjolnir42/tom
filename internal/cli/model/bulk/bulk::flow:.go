/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/cli/model/bulk"

import (
	"github.com/mjolnir42/tom/internal/cli/cmpl"
	"github.com/mjolnir42/tom/internal/cli/help"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func registerBulkFlow(app cli.App, run ActionFunc) *cli.App {
	handlerMap[proto.EntityFlow+`:`+proto.ActionList] = run(cmdBulkFlowList)
	handlerMap[proto.EntityFlow+`:`+proto.ActionShow] = run(cmdBulkFlowShow)

	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `flow`,
				Usage:       `Commands for maintaining flow collections`,
				Description: help.Text(proto.CmdFlow),
				Flags:       []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new flow collection`,
						Description:  help.Text(proto.CmdFlowAdd),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdBulkFlowAdd),
						BashComplete: cmpl.FlowAdd,
					},
					{
						Name:         `ensure`,
						Usage:        `Ensure a flow collection exists`,
						Description:  help.Text(proto.CmdFlowEnsure),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdBulkFlowEnsure),
						BashComplete: cmpl.FlowEnsure,
					},
					{
						Name:         `list`,
						Usage:        `List all existing flow collections`,
						Description:  help.Text(proto.CmdFlowList),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdBulkFlowList),
						BashComplete: cmpl.FlowList,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a flow collection`,
						Description:  help.Text(proto.CmdFlowShow),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdBulkFlowShow),
						BashComplete: cmpl.FlowShow,
					},
					{
						Name:        `property`,
						Usage:       `Commands for flow collection properties`,
						Description: help.Text(proto.CmdFlow),
						Flags:       []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Subcommands: []*cli.Command{
							{
								Name:         `set`,
								Usage:        `Set all properties of a flow collection`,
								Description:  help.Text(proto.CmdFlowPropSet),
								Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
								Action:       run(cmdBulkFlowPropSet),
								BashComplete: cmpl.FlowPropSet,
							},
							{
								Name:         `update`,
								Usage:        `Update properties of a flow collection`,
								Description:  help.Text(proto.CmdFlowPropUpdate),
								Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
								Action:       run(cmdBulkFlowPropUpdate),
								BashComplete: cmpl.FlowPropUpdate,
							},
							{
								Name:         `remove`,
								Usage:        `Remove properties of a flow collection`,
								Description:  help.Text(proto.CmdFlowPropRemove),
								Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
								Action:       run(cmdBulkFlowPropRemove),
								BashComplete: cmpl.FlowPropRemove,
							},
						},
					},
					{
						Name:         `remove`,
						Usage:        `Remove a flow collection`,
						Description:  help.Text(proto.CmdFlowRemove),
						Flags:        []cli.Flag{&cli.BoolFlag{Name: `verbose`, Aliases: []string{`v`}, Value: false, Hidden: true}},
						Action:       run(cmdBulkFlowRemove),
						BashComplete: cmpl.FlowRemove,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
