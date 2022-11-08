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
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ed25519"
)

func LoadCredentials(path string, phrase *string, lm *lhm.LogHandleMap, priv *epk.EncryptedPrivateKey, pub *ed25519.PublicKey, ctx *cli.Context) (bool, error) {
	var (
		err        error
		initialize bool = true
		rawPass    []byte
		fd         *os.File
	)

	// fix up credential path
	if path, err = filepath.Abs(path); err != nil {
		return false, err
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return false, err
	}
	if _, err = os.Open(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			lm.GetLogger(`application`).Infoln(`Credential directory missing, attempting create.`)

			if err = os.Mkdir(path, os.FileMode(0750)); err == nil {
				lm.GetLogger(`application`).Infoln(`successfully created credential directory.`)
			}
		} else {
			return false, err
		}
	}

	// load passphrase for private key from file
	if rawPass, err = ioutil.ReadFile(
		filepath.Join(path, `passphrase`),
	); err != nil && errors.Is(err, os.ErrNotExist) {
		if initialize {
			lm.GetLogger(`application`).Infoln(`initializing passphrase file`)
			if err = createPassphraseFile(filepath.Join(path, `passphrase`), phrase); err != nil {
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
		*phrase = string(mask(rawPass))
		lm.GetLogger(`application`).Infoln(`successfully loaded passphrase from file`)
	}

	// load keypair from file
	if fd, err = os.Open(
		filepath.Join(path, `machinekey.epk`),
	); err != nil && errors.Is(err, os.ErrNotExist) {
		if initialize {
			lm.GetLogger(`application`).Infoln(`initializing machine keypair`)
			if priv, pub, err = createKeypairFiles(
				path,
				*phrase,
			); err != nil {
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
		if priv, err = epk.ReadFrom(fd); err != nil {
			return false, err
		}
		lm.GetLogger(`application`).Infoln(`successfully loaded private key from file`)
	}

	// test loaded credentials
	if *pub, err = priv.Public(*phrase); err != nil {
		return false, err
	}
	lm.GetLogger(`application`).Infoln(`successfully unlocked public key from private key`)

	return initialize, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
