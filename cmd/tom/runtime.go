/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/urfave/cli/v2"
)

var client *resty.Client

func initCommon(c *cli.Context) error {

	cfg, err := configSetup(c)
	if err != nil {
		return err
	}

	client = resty.New().
		SetDisableWarn(true).
		SetHeader(`User-Agent`, fmt.Sprintf("%s %s", c.App.Name, c.App.Version)).
		SetHostURL(cfg.Run.API.String())

	// configure adm client library
	adm.ConfigureClient(client)
	adm.ConfigureJSONPostProcessor(cfg.ProcJSON)

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
