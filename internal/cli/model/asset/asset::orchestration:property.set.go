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
	proto.AssertCommandIsDefined(proto.CmdOrchestrationPropSet)
}

func cmdAssetOrchestrationPropSet(c *cli.Context) error {
	opts := map[string][]string{}
	props := []proto.PropertyDetail{}
	variable, once, required := adm.ArgumentsForCommand(proto.CmdOrchestrationPropSet)
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

	req := proto.NewOrchestrationRequest()
	req.Orchestration.Name = c.Args().First()
	req.Orchestration.Namespace = opts[`namespace`][0]
	req.Orchestration.Property = make(map[string]proto.PropertyDetail)

	if err := proto.ValidNamespace(req.Orchestration.Namespace); err != nil {
		return err
	}

	if err := proto.OnlyUnreserved(req.Orchestration.Name); err != nil {
		return err
	}

	for _, prop := range props {
		if err := proto.OnlyUnreserved(prop.Attribute); err != nil {
			return err
		}

		if err := proto.CheckPropertyConstraints(&prop); err != nil {
			return err
		}

		req.Orchestration.Property[prop.Attribute] = prop
	}

	spec := adm.Specification{
		Name: proto.CmdOrchestrationPropSet,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Orchestration.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
