/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	"encoding/base64"
	"flag"
	"os"
	"time"

	"github.com/Showmax/go-fqdn"
	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ed25519"
)

func registerMachineEnrollment(pub *ed25519.PublicKey, priv *epk.EncryptedPrivateKey, phrase string, ctx *cli.Context) error {
	var err error
	var msgBytes, sig []byte

	req := proto.NewUserRequest()
	req.User.LibraryName = `engineroom`
	if req.User.UserName, err = GetHash(pub); err != nil {
		return err
	}
	if req.User.FirstName, err = os.Hostname(); err != nil {
		return err
	}
	if req.User.LastName, err = fqdn.FqdnHostname(); err != nil {
		return err
	}
	req.User.Credential.Category = proto.CredentialPubKey
	req.User.Credential.Value = base64.StdEncoding.EncodeToString(*pub)

	req.Auth.Timestamp = time.Now().Format(time.RFC3339)
	req.Auth.UserID = req.User.UserName
	req.Auth.Fingerprint = req.User.UserName
	if err = req.CalculateDataHash(); err != nil {
		return err
	}
	if msgBytes, err = base64.StdEncoding.DecodeString(req.Auth.Sig.DataHash); err != nil {
		return err
	}
	if sig, err = priv.Sign(phrase, msgBytes); err != nil {
		return err
	}
	req.Auth.Sig.Signature = base64.StdEncoding.EncodeToString(sig)

	spec := adm.Specification{
		Name: proto.CmdMachEnrol,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.User.FormatMachineDNS(),
		},
		Body: req,
	}
	if ctx == nil {
		ctx = cli.NewContext(cli.NewApp(), flag.NewFlagSet(``, flag.ContinueOnError), nil)
	}

	return adm.Perform(spec, ctx)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
