/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

//go:generate go run ../../script/render_markdown.go ../../docs/tom/cmd_ref ../../internal/cli/help/rendered
//go:generate go-bindata -pkg help -ignore .gitignore -o ../../internal/cli/help/bindata.go -prefix "../../internal/cli/help/rendered/" ../../internal/cli/help/rendered/...

// global variables
var (
	// populated via Makefile
	tomVersion string
)

func main() {
	cli.CommandHelpTemplate = `{{.Description}}`

	app := cli.NewApp()
	app.Name = `tom`
	app.Usage = `Tom's Administrative Interface`
	app.Version = tomVersion
	app.EnableBashCompletion = true

	app = registerCommands(*app)
	//app = registerFlags(*app)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
