/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/tom/internal/config"
)

func createKeypairFiles(cfg *config.AuthConfiguration) error {
	var err error
	var fd *os.File

	if cfg.PrivEPK, cfg.PubKey, err = epk.New(cfg.Passphrase); err != nil {
		return err
	}

	if fd, err = os.Create(filepath.Join(cfg.CredPath, `machinekey.epk`)); err != nil {
		return err
	}
	if err = cfg.PrivEPK.Store(fd); err != nil {
		return err
	}
	fd.Close()

	if fd, err = os.Create(filepath.Join(cfg.CredPath, `machinekey.pub`)); err != nil {
		return err
	}
	fmt.Fprintf(fd, "%s %s\n", `ed25519-epk-pub`, base64.StdEncoding.EncodeToString(cfg.PubKey))
	fd.Close()
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
