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
	proto.AssertCommandIsDefined(proto.CmdContainerLink)
}

func cmdAssetContainerLink(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdContainerLink,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewContainerRequest()
	req.Container.TomID = c.Args().First()
	if err := req.Container.ParseTomID(); err != nil {
		return err
	}
	if err := proto.ValidNamespace(req.Container.Namespace); err != nil {
		return err
	}
	if err := proto.OnlyUnreserved(req.Container.Name); err != nil {
		return err
	}

	req.Container.Property = make(map[string]proto.PropertyDetail)
	for _, tgt := range opts[`is-equal`] {
		target := proto.Container{
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
		req.Container.Property[target.FormatDNS()] = proto.PropertyDetail{
			Attribute:  proto.ActionLink,
			Value:      target.FormatDNS(),
			ValidSince: `perpetual`,
			ValidUntil: `perpetual`,
		}
	}

	spec := adm.Specification{
		Name: proto.CmdContainerLink,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Container.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
