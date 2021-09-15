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
	proto.AssertCommandIsDefined(proto.CmdRuntimePropSet)
}

func cmdAssetRuntimePropSet(c *cli.Context) error {
	opts := map[string][]string{}
	props := []proto.PropertyDetail{}
	variable, once, required := adm.ArgumentsForCommand(proto.CmdRuntimePropSet)
	if err := adm.ParsePropertyArguments(
		opts,
		&props,
		c.Args().Tail(),
		variable,
		once,
		required,
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

	for _, prop := range props {
		if err := proto.OnlyUnreserved(prop.Attribute); err != nil {
			return err
		}

		if err := proto.CheckPropertyConstraints(&prop); err != nil {
			return err
		}

		req.Runtime.Property[prop.Attribute] = prop
	}

	spec := adm.Specification{
		Name: proto.CmdRuntimePropSet,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Runtime.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
