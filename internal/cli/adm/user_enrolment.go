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
	"time"

	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/cred"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

func RegisterUserEnrolment(cfg *config.AuthConfiguration, ctx *cli.Context) error {
	var err error

	if ctx == nil {
		ctx = cli.NewContext(cli.NewApp(), flag.NewFlagSet(``, flag.ContinueOnError), nil)
	}

	tx := time.Now().UTC()

	req := proto.NewUserRequest()
	req.User.LibraryName = cfg.IDLibrary
	req.User.UserName = cfg.UserName
	req.User.Credential.Category = proto.CredentialPubKey
	req.User.Credential.Value = base64.StdEncoding.EncodeToString(cfg.PubKey)

	req.Auth.Timestamp = tx.Format(time.RFC3339)
	req.Auth.UserID = req.User.UserName
	req.Auth.Fingerprint = cfg.Fingerprint

	req.Auth.CSR = &proto.DataCSR{
		UserID:       req.User.UserName,
		Library:      req.User.LibraryName,
		PublicKey:    req.User.Credential.Value,
		EnrolmentKey: enrolment, // package var, set via ConfigureEnrolmentKey()
		ValidFrom:    tx.Add(-1 * time.Second).Format(time.RFC3339),
		ValidUntil:   tx.Add(60 * time.Second).Format(time.RFC3339),
	}
	// CSR hash is included in the request body signature
	if err = req.Auth.CSR.CalculateDataHash(); err != nil {
		return err
	}

	if err = cred.SignRequestBody(&req, cfg); err != nil {
		return err
	}

	spec := Specification{
		Name: proto.CmdUserEnrol,
		Placeholder: map[string]string{
			proto.PlHoldTomID: req.User.FormatDNS(),
		},
		Body: req,
	}

	return Perform(spec, ctx)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
