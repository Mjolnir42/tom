/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	//"fmt"
	//"os"

	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func cmdMetaNamespaceAdd(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdNamespaceAdd,
		c,
		&opts,
	); err != nil {
		return err
	}

	req := proto.NewNamespaceRequest()
	req.Namespace.Property = make(map[string]proto.PropertyDetail)
	// set the namespace name
	req.Namespace.Property[`dict_name`] = proto.PropertyDetail{
		Attribute: `dict_name`,
		Value:     c.Args().First(),
	}

	// mandatory at-most-once argument
	req.Namespace.Property[`dict_type`] = proto.PropertyDetail{
		Attribute: `dict_type`,
		Value:     opts[`type`][0],
	}
	// optional at-most-once argument
	if _, ok := opts[`lookup-uri`]; ok {
		req.Namespace.Property[`dict_uri`] = proto.PropertyDetail{
			Attribute: `dict_uri`,
			Value:     opts[`lookup-uri`][0],
		}
	}
	// optional at-most-once argument
	if _, ok := opts[`lookup-key`]; ok {
		req.Namespace.Property[`dict_lookup`] = proto.PropertyDetail{
			Attribute: `dict_lookup`,
			Value:     opts[`lookup-key`][0],
		}
	}
	// optional at-most-once argument
	if _, ok := opts[`entities`]; ok {
		req.Namespace.Property[`dict_ntt_list`] = proto.PropertyDetail{
			Attribute: `dict_ntt_list`,
			Value:     opts[`entities`][0],
		}
	}

	// optional arguments that can each be given multiple times
	req.Namespace.Attributes = make([]proto.AttributeDefinition, 0)
	if _, ok := opts[`std-attr`]; ok {
		for _, std := range opts[`std-attr`] {
			req.Namespace.Attributes = append(
				req.Namespace.Attributes,
				proto.AttributeDefinition{Key: std, Unique: false},
			)
		}
	}
	if _, ok := opts[`uniq-attr`]; ok {
		for _, uniq := range opts[`uniq-attr`] {
			req.Namespace.Attributes = append(
				req.Namespace.Attributes,
				proto.AttributeDefinition{Key: uniq, Unique: false},
			)
		}
	}

	return adm.Perform(`postbody`, `/namespace/`, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
