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
	proto.AssertCommandIsDefined(proto.CmdRuntimeLink)
}

func cmdAssetRuntimeLink(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdRuntimeLink,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewRuntimeRequest()
	req.Runtime.TomID = c.Args().First()
	if err := req.Runtime.ParseTomID(); err != nil {
		return err
	}
	if err := proto.ValidNamespace(req.Runtime.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(req.Runtime.Name); err != nil {
		return err
	}

	req.Runtime.Property = make(map[string]proto.PropertyDetail)
	for _, tgt := range opts[`is-equal`] {
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
		req.Runtime.Property[proto.MetaPropertyCmdLink+`::`+target.FormatDNS()] = proto.PropertyDetail{
			Attribute:  proto.MetaPropertyCmdLink,
			Value:      target.FormatDNS(),
			ValidSince: `perpetual`,
			ValidUntil: `perpetual`,
		}
	}

	spec := adm.Specification{
		Name: proto.CmdRuntimeLink,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Runtime.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
