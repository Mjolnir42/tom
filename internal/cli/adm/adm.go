/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * Copyright (c) 2021, Jörg Pernfuß <joerg.pernfuss@ionos.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"github.com/go-resty/resty/v2"
	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/tom/internal/cli/db"
)

var (
	client        *resty.Client
	cache         *db.DB
	async         bool
	jobSave       bool
	postProcessor string
	idLibID       string
	userID        string
	authenticate  bool
	priv          *epk.EncryptedPrivateKey
)

func ConfigureClient(c *resty.Client) {
	client = c
}

func ConfigureCache(c *db.DB) {
	cache = c
}

func ConfigureIdentity(lib, user string) {
	idLibID = lib
	userID = user

	if priv != nil {
		authenticate = true
	}
}

func ConfigureEPK(pk *epk.EncryptedPrivateKey) {
	priv = pk

	if idLibID != `` && userID != `` {
		authenticate = true
	}
}

func ActivateAsyncWait(b bool) {
	async = b
}

func AutomaticJobSave(b bool) {
	jobSave = b
}

func ConfigureJSONPostProcessor(p string) {
	postProcessor = p
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
