/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
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
	proto.AssertCommandIsDefined(proto.CmdServerPropSet)
}

func cmdAssetServerPropSet(c *cli.Context) error {
	opts := map[string][]string{}
	props := []proto.PropertyDetail{}
	variable, once, required := adm.ArgumentsForCommand(proto.CmdServerPropSet)
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

	req := proto.NewServerRequest()
	req.Server.Name = c.Args().First()
	req.Server.Namespace = opts[`namespace`][0]
	req.Server.Property = make(map[string]proto.PropertyDetail)

	if err := proto.ValidNamespace(req.Server.Namespace); err != nil {
		return err
	}

	if err := proto.OnlyUnreserved(req.Server.Name); err != nil {
		return err
	}

	for _, prop := range props {
		if err := proto.OnlyUnreserved(prop.Attribute); err != nil {
			return err
		}

		if err := proto.CheckPropertyConstraints(&prop); err != nil {
			return err
		}

		req.Server.Property[prop.Attribute] = prop
	}

	spec := adm.Specification{
		Name: proto.CmdServerPropSet,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Server.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
