/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package core // import "github.com/mjolnir42/tom/internal/core/"

import (
	"github.com/mjolnir42/tom/internal/model/asset"
	"github.com/mjolnir42/tom/internal/model/bulk"
	"github.com/mjolnir42/tom/internal/model/iam"
	"github.com/mjolnir42/tom/internal/model/meta"
	"github.com/mjolnir42/tom/internal/model/supervisor"
)

// Start launches all application handlers
func (x *Core) Start() {
	// iam model
	iam.HandleRegister(x.hm, x.conf.QueueLen)

	// meta Model
	meta.HandleRegister(x.hm, x.conf.QueueLen)

	// asset model
	asset.HandleRegister(x.hm, x.conf.QueueLen)

	// bulk model
	bulk.HandleRegister(x.hm, x.conf.QueueLen)

	// supervisor engine
	supervisor.HandleRegister(x.hm, x.conf.QueueLen)

	for handlerName := range x.hm.Range() {
		x.hm.Configure(
			handlerName,
			x.db,
			x.lm,
		)
		// start the handler in a goroutine
		x.lm.GetLogger(`application`).Infof(
			"Core running handler: %s",
			handlerName,
		)
		x.hm.Run(handlerName)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
