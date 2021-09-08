/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/cli/model/meta"

import (
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespacePropRemove)
}

func cmdMetaNamespacePropRemove(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdNamespacePropRemove,
		c,
		&opts,
	); err != nil {
		return err
	}

	if err := proto.ValidNamespace(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewNamespaceRequest()
	req.Namespace.Name = c.Args().First()
	req.Namespace.Property = make(map[string]proto.PropertyDetail)

	for _, prop := range opts[`property`] {
		if err := proto.OnlyUnreserved(prop); err != nil {
			return err
		}

		req.Namespace.Property[prop] = proto.PropertyDetail{
			Attribute: prop,
		}
	}

	spec := adm.Specification{
		Name: proto.CmdNamespacePropRemove,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Namespace.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
