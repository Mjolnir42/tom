/*-
 * Copyright (c) 2021, Jörg Pernfuß
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
	proto.AssertCommandIsDefined(proto.CmdRuntimePropRemove)
}

func cmdAssetRuntimePropRemove(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdRuntimePropRemove,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewRuntimeRequest()
	req.Runtime.Name = c.Args().First()
	req.Runtime.Namespace = opts[`namespace`][0]
	req.Runtime.Property = make(map[string]proto.PropertyDetail)

	if err := proto.ValidNamespace(req.Runtime.Namespace); err != nil {
		return err
	}

	if err := proto.OnlyUnreserved(req.Runtime.Name); err != nil {
		return err
	}

	for _, prop := range opts[`property`] {
		if err := proto.OnlyUnreserved(prop); err != nil {
			return err
		}

		req.Runtime.Property[prop] = proto.PropertyDetail{
			Attribute: prop,
		}
	}

	spec := adm.Specification{
		Name: proto.CmdRuntimePropRemove,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Runtime.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
