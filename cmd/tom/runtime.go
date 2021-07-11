/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	//"fmt"
	//"os"

	"github.com/urfave/cli/v2"
)

func runtime(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		return action(c)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
