/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/cli/model/bulk"

import (
	"github.com/urfave/cli/v2"
)

type ActionFunc func(cli.ActionFunc) cli.ActionFunc
type Registry map[string]func(*cli.Context) error

var handlerMap Registry

func Register(app cli.App, run ActionFunc, r Registry) *cli.App {
	handlerMap = r

	app = *registerBulkFlow(app, run)
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
