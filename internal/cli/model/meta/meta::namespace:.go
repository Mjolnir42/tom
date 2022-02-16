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
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func registerMetaNamespace(app cli.App, run ActionFunc) *cli.App {
	handlerMap[proto.EntityNamespace+`:`+proto.ActionList] = run(cmdMetaNamespaceList)
	handlerMap[proto.EntityNamespace+`:`+proto.ActionShow] = run(cmdMetaNamespaceShow)

	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:        `namespace`,
				Usage:       `Commands for maintaining namespaces`,
				Description: help.Text(proto.CmdNamespace),
				Subcommands: []*cli.Command{
					{
						Name:         `add`,
						Usage:        `Create a new namespace`,
						Description:  help.Text(proto.CmdNamespaceAdd),
						Action:       run(cmdMetaNamespaceAdd),
						BashComplete: cmpl.NamespaceAdd,
					},
					{
						Name:         `list`,
						Usage:        `List all existing namespaces`,
						Description:  help.Text(proto.CmdNamespaceList),
						Action:       run(cmdMetaNamespaceList),
						BashComplete: cmpl.NamespaceList,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a namespaces`,
						Description:  help.Text(proto.CmdNamespaceShow),
						Action:       run(cmdMetaNamespaceShow),
						BashComplete: cmpl.NamespaceShow,
					},
					{
						Name:        `attribute`,
						Usage:       `Section for attribute commands`,
						Description: help.Text(proto.CmdNamespace),
						Subcommands: []*cli.Command{
							{
								Name:         `add`,
								Usage:        `Add new attributes to a namespace`,
								Description:  help.Text(proto.CmdNamespaceAttrAdd),
								Action:       run(cmdMetaNamespaceAttrAdd),
								BashComplete: cmpl.NamespaceAttrAdd,
							},
							{
								Name:         `remove`,
								Usage:        `Remove attributes from a namespace`,
								Description:  help.Text(proto.CmdNamespaceAttrRemove),
								Action:       run(cmdMetaNamespaceAttrRemove),
								BashComplete: cmpl.NamespaceAttrRemove,
							},
						},
					},
					{
						Name:        `property`,
						Usage:       `Section for property commands`,
						Description: help.Text(proto.CmdNamespace),
						Subcommands: []*cli.Command{
							{
								Name:         `set`,
								Usage:        `Set all properties of a namespace`,
								Description:  help.Text(proto.CmdNamespacePropSet),
								Action:       run(cmdMetaNamespacePropSet),
								BashComplete: cmpl.NamespacePropSet,
							},
							{
								Name:         `update`,
								Usage:        `Update properties of a namespace`,
								Description:  help.Text(proto.CmdNamespacePropUpdate),
								Action:       run(cmdMetaNamespacePropUpdate),
								BashComplete: cmpl.NamespacePropUpdate,
							},
							{
								Name:         `remove`,
								Usage:        `Remove properties of a namespace`,
								Description:  help.Text(proto.CmdNamespacePropRemove),
								Action:       run(cmdMetaNamespacePropRemove),
								BashComplete: cmpl.NamespacePropRemove,
							},
						},
					},
					{
						Name:         `remove`,
						Usage:        `Remove a namespace`,
						Description:  help.Text(proto.CmdNamespaceRemove),
						Action:       run(cmdMetaNamespaceRemove),
						BashComplete: cmpl.NamespaceRemove,
					},
				},
			},
		}...,
	)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
