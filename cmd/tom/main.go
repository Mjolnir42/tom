/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
)

//go:generate go run ../../script/render_markdown.go ../../docs/tom/cmd_ref ../../internal/cli/help/rendered
//go:generate cp ../../tools/zsh_autocomplete ../../internal/cli/help/rendered/zsh_autocomplete.fmt
//go:generate go-bindata -pkg help -ignore .gitignore -o ../../internal/cli/help/bindata.go -prefix "../../internal/cli/help/rendered/" ../../internal/cli/help/rendered/...

// global variables
var (
	// populated via Makefile
	tomVersion string
	// used for generating the RequestID for client errors
	namespaceTom = uuid.Must(uuid.FromString(`ffffffff-0000-5000-0000-ffffffffffff`))
	// setup by runtime function to provide cli arguments to error
	// output processing
	errorContext *cli.Context
)

func main() {
	cli.CommandHelpTemplate = `{{.Description}}`

	app := cli.NewApp()
	app.Name = `tom`
	app.Usage = `Tom's Administrative Interface`
	app.Version = tomVersion
	app.EnableBashCompletion = true

	app = registerCommands(*app)

	if err := app.Run(os.Args); err != nil {
		data, jsonError := json.Marshal(&proto.Result{
			StatusCode: 400,
			RequestID:  uuid.NewV5(namespaceTom, app.Name+tomVersion).String(),
			ErrorText:  err.Error(),
		})
		if jsonError != nil {
			// dual errors: bail on formatted error output
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, jsonError)
			os.Exit(2)
		}

		if formatError := adm.FormatOut(errorContext, data, proto.TemplateCommand); formatError != nil {
			// dual errors: bail on formatted error output
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, formatError)
			os.Exit(2)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
