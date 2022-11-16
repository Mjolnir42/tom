/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

func LoadCredentials(cfg *config.AuthConfiguration, lm *lhm.LogHandleMap) (bool, error) {
	var (
		err        error
		initialize bool = true
		rawPass    []byte
		fd         *os.File
	)

	// fix up credential path
	if cfg.CredPath, err = filepath.Abs(cfg.CredPath); err != nil {
		return false, err
	}
	if cfg.CredPath, err = filepath.EvalSymlinks(cfg.CredPath); err != nil {
		return false, err
	}
	if _, err = os.Open(cfg.CredPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			lm.GetLogger(`application`).Infoln(`Credential directory missing, attempting create.`)

			if err = os.Mkdir(cfg.CredPath, os.FileMode(0750)); err == nil {
				lm.GetLogger(`application`).Infoln(`successfully created credential directory.`)
			}
		} else {
			return false, err
		}
	}

	// load passphrase for private key from file
	if rawPass, err = ioutil.ReadFile(
		filepath.Join(cfg.CredPath, `passphrase`),
	); err != nil && errors.Is(err, os.ErrNotExist) {
		if initialize {
			lm.GetLogger(`application`).Infoln(`initializing passphrase file`)
			if err = createPassphraseFile(filepath.Join(cfg.CredPath, `passphrase`), cfg); err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	} else if err != nil {
		// error opening file
		return false, err
	} else {
		// successfully read passphrase file, deactivate initialize mode
		initialize = false
		if len(rawPass) == 0 {
			return false, errors.New(`passphrase file is empty`)
		}
		cfg.Passphrase = string(mask(rawPass))
		lm.GetLogger(`application`).Infoln(`successfully loaded passphrase from file`)
	}

	// load keypair from file
	if fd, err = os.Open(
		filepath.Join(cfg.CredPath, `machinekey.epk`),
	); err != nil && errors.Is(err, os.ErrNotExist) {
		if initialize {
			lm.GetLogger(`application`).Infoln(`initializing machine keypair`)
			if err = createKeypairFiles(cfg); err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	} else if err != nil {
		// error opening file
		return false, err
	} else {
		// successfully read private keyfile, deactivate initialize mode
		initialize = false
		if cfg.PrivEPK, err = epk.ReadFrom(fd); err != nil {
			return false, err
		}
		lm.GetLogger(`application`).Infoln(`successfully loaded private key from file`)
	}

	// test loaded credentials
	if cfg.PubKey, err = cfg.PrivEPK.Public(cfg.Passphrase); err != nil {
		return false, err
	}
	lm.GetLogger(`application`).Infoln(`successfully unlocked public key from private key`)
	// pre-calculate publickey fingerprint hash
	if cfg.Fingerprint, err = GetHash(cfg.PubKey); err != nil {
		return false, err
	}

	return initialize, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
