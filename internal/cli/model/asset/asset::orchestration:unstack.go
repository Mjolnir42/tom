/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/cli/model/asset"

import (
	"fmt"

	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdOrchestrationUnstack)
}

func cmdAssetOrchestrationUnstack(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdOrchestrationUnstack,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewOrchestrationRequest()
	if _, ok := opts[`namespace`]; ok {
		req.Orchestration.Namespace = opts[`namespace`][0]
		req.Orchestration.Name = c.Args().First()
	} else {
		req.Orchestration.TomID = c.Args().First()
		if err := req.Orchestration.ParseTomID(); err != nil {
			return err
		}
	}

	if err := proto.ValidNamespace(req.Orchestration.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(req.Orchestration.Name); err != nil {
		return err
	}

	req.Orchestration.Property = make(map[string]proto.PropertyDetail)

	for i := range opts[`unprovide`] {
		prop := proto.PropertyDetail{
			Attribute: proto.MetaPropertyCmdUnstack,
		}
		if err, entityType, ntt := proto.ParseTomID(
			opts[`unprovide`][i],
		); err != nil {
			return err
		} else {
			prop.Value = ntt.FormatDNS()

			if err := proto.ValidNamespace(ntt.ExportNamespace()); err != nil {
				return err
			}
			if err := proto.OnlyUnreserved(ntt.ExportName()); err != nil {
				return err
			}

			switch entityType {
			case proto.EntityRuntime:
			default:
				return fmt.Errorf("Invalid stacking target class: %s", entityType)
			}
		}
		if _, ok := opts[`since`]; ok {
			prop.ValidSince = opts[`since`][0]
		}

		req.Orchestration.Property[proto.MetaPropertyCmdUnstack+`::`+prop.Value] = prop
	}

	spec := adm.Specification{
		Name: proto.CmdOrchestrationUnstack,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Orchestration.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
