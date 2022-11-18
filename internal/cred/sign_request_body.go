/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	"encoding/base64"

	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/pkg/proto"
)

func SignRequestBody(req *proto.Request, cfg *config.AuthConfiguration) error {
	var err error
	var msgBytes, sig []byte

	if err = req.CalculateDataHash(); err != nil {
		return err
	}
	if msgBytes, err = base64.StdEncoding.DecodeString(req.Auth.Sig.DataHash); err != nil {
		return err
	}
	if sig, err = cfg.PrivEPK.Sign(cfg.Passphrase, msgBytes); err != nil {
		return err
	}
	req.Auth.Sig.Signature = base64.StdEncoding.EncodeToString(sig)

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
