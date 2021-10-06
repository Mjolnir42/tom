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
	proto.AssertCommandIsDefined(proto.CmdContainerAdd)
}

func cmdAssetContainerAdd(c *cli.Context) error {
	opts := map[string][]string{}
	props := []proto.PropertyDetail{}
	variable, once, required := adm.ArgumentsForCommand(proto.CmdContainerAdd)
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

	req := proto.NewContainerRequest()
	if err := proto.ValidNamespace(opts[`namespace`][0]); err != nil {
		return err
	}
	req.Container.Namespace = opts[`namespace`][0]
	req.Container.Property = make(map[string]proto.PropertyDetail)
	for _, prop := range props {
		if err := proto.OnlyUnreserved(prop.Attribute); err != nil {
			return err
		}
		req.Container.Property[prop.Attribute] = prop
	}

	if err := proto.OnlyUnreserved(c.Args().First()); err != nil {
		return err
	}
	req.Container.Property[`name`] = proto.PropertyDetail{
		Attribute: `name`,
		Value:     c.Args().First(),
	}

	if _, ok := opts[`since`]; ok {
		prop := req.Container.Property[`name`]
		prop.ValidSince = opts[`since`][0]
		req.Container.Property[`name`] = prop
	}

	if _, ok := opts[`until`]; ok {
		prop := req.Container.Property[`name`]
		prop.ValidUntil = opts[`until`][0]
		req.Container.Property[`name`] = prop
	}

	req.Container.Property[`type`] = proto.PropertyDetail{
		Attribute:  `type`,
		Value:      opts[`type`][0],
		ValidSince: `perpetual`,
		ValidUntil: `perpetual`,
	}

	spec := adm.Specification{
		Name: proto.CmdContainerAdd,
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
