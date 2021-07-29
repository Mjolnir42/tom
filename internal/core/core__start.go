/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package core // import "github.com/mjolnir42/tom/internal/core/"

import (
	"github.com/mjolnir42/tom/internal/model/asset"
	"github.com/mjolnir42/tom/internal/model/iam"
	"github.com/mjolnir42/tom/internal/model/meta"
)

// Start launches all application handlers
func (x *Core) Start() {
	// IAM Model
	x.hm.Add(iam.NewLibraryReadHandler(x.conf.QueueLen))
	x.hm.Add(iam.NewLibraryWriteHandler(x.conf.QueueLen))
	x.hm.Add(iam.NewUserReadHandler(x.conf.QueueLen))
	x.hm.Add(iam.NewUserWriteHandler(x.conf.QueueLen))

	x.hm.Add(meta.NewNamespaceReadHandler(x.conf.QueueLen))
	x.hm.Add(meta.NewNamespaceWriteHandler(x.conf.QueueLen))

	x.hm.Add(asset.NewServerReadHandler(x.conf.QueueLen))
	x.hm.Add(asset.NewServerWriteHandler(x.conf.QueueLen))
	x.hm.Add(asset.NewRuntimeReadHandler(x.conf.QueueLen))
	x.hm.Add(asset.NewRuntimeWriteHandler(x.conf.QueueLen))

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
