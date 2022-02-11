/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/cli/model/asset"

import (
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdServerUnstack)
}

func cmdAssetServerUnstack(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdServerUnstack,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewServerRequest()
	if _, ok := opts[`namespace`]; ok {
		req.Server.Namespace = opts[`namespace`][0]
		req.Server.Name = c.Args().First()
	} else {
		req.Server.TomID = c.Args().First()
		if err := req.Server.ParseTomID(); err != nil {
			return err
		}
	}

	if err := proto.ValidNamespace(req.Server.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(req.Server.Name); err != nil {
		return err
	}

	spec := adm.Specification{
		Name: proto.CmdServerUnstack,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Server.FormatDNS(),
		},
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
