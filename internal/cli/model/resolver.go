/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package model // import "github.com/mjolnir42/tom/internal/cli/model"

import (
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/internal/cli/cmpl"
	"github.com/mjolnir42/tom/internal/cli/help"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func registerResolver(app cli.App, run ActionFunc) *cli.App {
	app.Commands = append(app.Commands,
		[]*cli.Command{
			{
				Name:         `query`,
				Usage:        `Query TOM by tomID`,
				Description:  help.Text(`model::query:resolver`),
				Action:       cmdQueryResolver,
				BashComplete: cmpl.QueryResolver,
			},
		}...,
	)
	return &app
}

func cmdQueryResolver(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	switch {
	case proto.IsTomID(c.Args().First()):
		_, ntt, _ := proto.ParseTomID(c.Args().First())
		return handlerMap[ntt+`:`+proto.ActionShow](c)
	case proto.IsWildcardTomID(c.Args().First()):
		_, _, ntt := proto.ParseTomIDWildcard(c.Args().First())
		return handlerMap[ntt+`:`+proto.ActionList](c)
	default:
		return proto.ErrParsingTomID
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
