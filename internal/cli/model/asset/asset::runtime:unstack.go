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
	proto.AssertCommandIsDefined(proto.CmdRuntimeUnstack)
}

func cmdAssetRuntimeUnstack(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdRuntimeUnstack,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewRuntimeRequest()
	if _, ok := opts[`namespace`]; ok {
		req.Runtime.Namespace = opts[`namespace`][0]
		req.Runtime.Name = c.Args().First()
	} else {
		req.Runtime.TomID = c.Args().First()
		if err := req.Runtime.ParseTomID(); err != nil {
			return err
		}
	}

	if err := proto.ValidNamespace(req.Runtime.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(req.Runtime.Name); err != nil {
		return err
	}

	spec := adm.Specification{
		Name: proto.CmdRuntimeUnstack,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Runtime.FormatDNS(),
		},
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
