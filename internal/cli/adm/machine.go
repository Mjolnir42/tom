/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm // import "github.com/mjolnir42/tom/internal/cli/adm"

import (
	"encoding/base64"
	"flag"
	"os"
	"time"

	"github.com/Showmax/go-fqdn"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/cred"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func RegisterMachineEnrollment(cfg *config.AuthConfiguration, ctx *cli.Context) error {
	var err error

	ConfigureEPK(cfg.PrivEPK, cfg.Passphrase)

	req := proto.NewUserRequest()
	req.User.LibraryName = `engineroom`
	if req.User.UserName, err = cred.GetHash(cfg.PubKey); err != nil {
		return err
	}
	if req.User.FirstName, err = os.Hostname(); err != nil {
		return err
	}
	if req.User.LastName, err = fqdn.FqdnHostname(); err != nil {
		return err
	}
	req.User.Credential.Category = proto.CredentialPubKey
	req.User.Credential.Value = base64.StdEncoding.EncodeToString(cfg.PubKey)

	req.Auth.Timestamp = time.Now().Format(time.RFC3339)
	req.Auth.UserID = req.User.UserName
	req.Auth.Fingerprint = req.User.UserName

	if err = cred.SignRequestBody(&req, cfg); err != nil {
		return err
	}

	spec := Specification{
		Name: proto.CmdMachEnrol,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.User.FormatMachineDNS(),
		},
		Body: req,
	}
	if ctx == nil {
		ctx = cli.NewContext(cli.NewApp(), flag.NewFlagSet(``, flag.ContinueOnError), nil)
	}

	return Perform(spec, ctx)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
