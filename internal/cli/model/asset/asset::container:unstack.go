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
	proto.AssertCommandIsDefined(proto.CmdContainerUnstack)
}

func cmdAssetContainerUnstack(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdContainerUnstack,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewContainerRequest()
	if _, ok := opts[`namespace`]; ok {
		req.Container.Namespace = opts[`namespace`][0]
		req.Container.Name = c.Args().First()
	} else {
		req.Container.TomID = c.Args().First()
		if err := req.Container.ParseTomID(); err != nil {
			return err
		}
	}

	if err := proto.ValidNamespace(req.Container.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(req.Container.Name); err != nil {
		return err
	}

	spec := adm.Specification{
		Name: proto.CmdContainerUnstack,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Container.FormatDNS(),
		},
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
