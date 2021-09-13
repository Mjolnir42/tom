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
	proto.AssertCommandIsDefined(proto.CmdRuntimeList)
}

func cmdAssetRuntimeList(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.VariadicDirect(
		proto.CmdRuntimeList,
		c,
		&opts,
	); err != nil {
		return err
	}

	spec := adm.Specification{
		Name: proto.CmdRuntimeList,
	}
	if _, ok := opts[`namespace`]; ok {
		if err := proto.ValidNamespace(opts[`namespace`][0]); err != nil {
			return err
		}

		spec.QueryParams = &map[string]string{
			`namespace`: opts[`namespace`][0],
		}
	}
	return adm.Perform(spec, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
