/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/cli/model/meta"

import (
	"fmt"

	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespaceAttrRemove)
}

func cmdMetaNamespaceAttrRemove(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdNamespaceAttrRemove,
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

	// optional arguments that can each be given multiple times
	req.Namespace.Attributes = make([]proto.AttributeDefinition, 0)
	if _, ok := opts[`std-attr`]; ok {
		for _, std := range opts[`std-attr`] {
			if err := proto.OnlyUnreserved(std); err != nil {
				return err
			}

			req.Namespace.Attributes = append(
				req.Namespace.Attributes,
				proto.AttributeDefinition{Key: std, Unique: false},
			)
		}
	}
	if _, ok := opts[`uniq-attr`]; ok {
		for _, uniq := range opts[`uniq-attr`] {
			if err := proto.OnlyUnreserved(uniq); err != nil {
				return err
			}

			req.Namespace.Attributes = append(
				req.Namespace.Attributes,
				proto.AttributeDefinition{Key: uniq, Unique: true},
			)
		}
	}

	// check that at least one attribute was defined
	if len(req.Namespace.Attributes) == 0 {
		return fmt.Errorf("Specified no attributes to add")
	}

	spec := adm.Specification{
		Name: proto.CmdNamespaceAttrRemove,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.Namespace.FormatDNS(),
		},
		Body: req,
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
