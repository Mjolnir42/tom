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

func initCommon(c *cli.Context) error {

	cfg, err := configSetup(c)
	if err != nil {
		return err
	}

	return nil
}

func runtime(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {

		if err := initCommon(c); err != nil {
			return err
		}
		return action(c)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
