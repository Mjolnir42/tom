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
	proto.AssertCommandIsDefined(proto.CmdContainerStack)
}

func cmdAssetContainerStack(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdContainerStack,
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

	req.Container.Property = make(map[string]proto.PropertyDetail)
	tgt := opts[`provided-by`][0]
	target := proto.Runtime{
		TomID: tgt,
	}
	if err := target.ParseTomID(); err != nil {
		return err
	}
	if err := proto.ValidNamespace(target.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(target.Name); err != nil {
		return err
	}
	prop := proto.PropertyDetail{
		Attribute: proto.MetaPropertyCmdStack,
		Value:     target.FormatDNS(),
	}

	if _, ok := opts[`since`]; ok {
		prop.ValidSince = opts[`since`][0]
	}

	if _, ok := opts[`until`]; ok {
		prop.ValidUntil = opts[`until`][0]
	}

	req.Container.Property[proto.MetaPropertyCmdUnstack+`::`+prop.Value] = prop

	spec := adm.Specification{
		Name: proto.CmdContainerStack,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Container.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
