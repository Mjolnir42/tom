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
	proto.AssertCommandIsDefined(proto.CmdServerShow)
}

func cmdAssetServerShow(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicArguments(
		proto.CmdServerShow,
		c,
		&opts,
	); err != nil {
		return err
	}

	r := proto.Server{}
	if _, ok := opts[`namespace`]; ok {
		r.Namespace = opts[`namespace`][0]
		r.Name = c.Args().First()
		if err := proto.ValidNamespace(r.Namespace); err != nil {
			return err
		}
		if err := proto.OnlyUnreserved(r.Name); err != nil {
			return err
		}
	} else {
		r.TomID = c.Args().First()
		if err := r.ParseTomID(); err != nil {
			return err
		}
	}

	spec := adm.Specification{
		Name: proto.CmdServerShow,
		Placeholder: map[string]string{
			proto.PlHoldTomID: r.FormatDNS(),
		},
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
