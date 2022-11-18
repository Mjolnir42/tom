/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"net/url"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/cred"
	"github.com/urfave/cli/v2"
)

func configSetup(c *cli.Context) (*config.ClientConfig, bool, error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, false, err
	}

	var confPath string
	switch {
	case c.IsSet(`config`) && filepath.IsAbs(c.String(`config`)):
		confPath = filepath.Clean(c.String(`config`))
	case c.IsSet(`config`):
		confPath = filepath.Clean(filepath.Join(home, ".tom", c.String(`config`)))
	default:
		confPath = filepath.Clean(filepath.Join(home, ".tom", "tom.conf"))
	}

	cfg := &config.ClientConfig{Run: config.RunTimeConfig{}, Auth: &config.AuthConfiguration{}}
	if err = cfg.PopulateFromFile(confPath); err != nil {
		return nil, false, err
	}
	cfg.Run.API, err = url.Parse(cfg.API)
	if err != nil {
		return nil, false, err
	}

	cfg.Run.PathLogs = filepath.Clean(cfg.LogDir)
	if cfg.CAFile != `` {
		cfg.Run.PathCA = filepath.Clean(cfg.CAFile)
	}

	switch cfg.Auth.CredPath {
	case ``:
		cfg.Auth.CredPath = filepath.Clean(filepath.Join(home, ".tom"))
	default:
		cfg.Auth.CredPath = filepath.Clean(cfg.Auth.CredPath)
	}

	var initialize bool
	if initialize, err = cred.LoadCredentials(
		cfg.Auth,
		nil,
	); err != nil {
		return nil, false, err
	}

	switch cfg.Auth.UseFingerprint {
	case true:
		adm.ConfigureIdentity(cfg.Auth.IDLibrary, cfg.Auth.Fingerprint)
	default:
		adm.ConfigureIdentity(cfg.Auth.IDLibrary, cfg.Auth.UserName)
	}

	adm.ConfigureEPK(cfg.Auth.PrivEPK, cfg.Auth.Passphrase)

	return cfg, initialize, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
