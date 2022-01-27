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
	proto.AssertCommandIsDefined(proto.CmdRuntimeStack)
}

func cmdAssetRuntimeStack(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdRuntimeStack,
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

	req.Runtime.Property = make(map[string]proto.PropertyDetail)
	prop := proto.PropertyDetail{
		Attribute: proto.MetaPropertyCmdStack,
	}

	if err, entityType, ntt := proto.ParseTomID(
		opts[`runs-on`][0],
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
		case proto.EntityRuntime, proto.EntityServer, proto.EntityOrchestration:
		default:
			return fmt.Errorf("Invalid stacking target class: %s", entityType)
		}
	}

	if _, ok := opts[`since`]; ok {
		prop.ValidSince = opts[`since`][0]
	}

	if _, ok := opts[`until`]; ok {
		prop.ValidUntil = opts[`until`][0]
	}

	req.Runtime.Property[proto.MetaPropertyCmdStack+`::`+prop.Value] = prop

	spec := adm.Specification{
		Name: proto.CmdRuntimeStack,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Runtime.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
