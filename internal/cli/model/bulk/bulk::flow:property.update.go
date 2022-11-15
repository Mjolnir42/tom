/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/cli/model/bulk"

import (
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdFlowPropUpdate)
}

func cmdBulkFlowPropUpdate(c *cli.Context) error {
	opts := map[string][]string{}
	props := []proto.PropertyDetail{}
	variable, once, required := adm.ArgumentsForCommand(proto.CmdFlowPropUpdate)
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

	req := proto.NewFlowRequest()
	req.Flow.Name = c.Args().First()
	req.Flow.Namespace = opts[`namespace`][0]
	req.Flow.Property = make(map[string]proto.PropertyDetail)

	if err := proto.ValidNamespace(req.Flow.Namespace); err != nil {
		return err
	}

	if err := proto.OnlyUnreserved(req.Flow.Name); err != nil {
		return err
	}

	for _, prop := range props {
		if err := proto.OnlyUnreserved(prop.Attribute); err != nil {
			return err
		}

		if err := proto.CheckPropertyConstraints(&prop); err != nil {
			return err
		}

		req.Flow.Property[prop.Attribute] = prop
	}

	spec := adm.Specification{
		Name: proto.CmdFlowPropUpdate,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Flow.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
