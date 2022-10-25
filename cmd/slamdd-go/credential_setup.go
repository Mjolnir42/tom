/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/tom/cmd/slamdd-go"

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/bits"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjolnir42/epk"
)

func loadCredentials() error {
	var (
		path       string
		err        error
		initialize bool = true
		rawPass    []byte
		fd         *os.File
	)

	// fix up credential path
	if path, err = filepath.Abs(SlamCfg.CredPath); err != nil {
		return err
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return err
	}
	if _, err = os.Open(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			lm.GetLogger(`application`).Infoln(`Credential directory missing, attempting create.`)

			if err = os.Mkdir(path, os.FileMode(0750)); err == nil {
				lm.GetLogger(`application`).Infoln(`successfully created credential directory.`)
			}
		} else {
			return err
		}
	}

	// load passphrase for private key from file
	if rawPass, err = ioutil.ReadFile(
		filepath.Join(path, `passphrase`),
	); err != nil && errors.Is(err, os.ErrNotExist) {
		if initialize {
			lm.GetLogger(`application`).Infoln(`initializing passphrase file`)
			if err = createPassphraseFile(filepath.Join(path, `passphrase`)); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if err != nil {
		// error opening file
		return err
	} else {
		// successfully read passphrase file, deactivate initialize mode
		initialize = false
		if len(rawPass) == 0 {
			return errors.New(`passphrase file is empty`)
		}
		SlamCfg.Passphrase = string(mask(rawPass))
		lm.GetLogger(`application`).Infoln(`successfully loaded passphrase from file`)
	}

	// load keypair from file
	if fd, err = os.Open(
		filepath.Join(path, `machinekey.epk`),
	); err != nil && errors.Is(err, os.ErrNotExist) {
		if initialize {
			lm.GetLogger(`application`).Infoln(`initializing machine keypair`)
			if err = createKeypairFiles(path); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if err != nil {
		// error opening file
		return err
	} else {
		// successfully read private keyfile, deactivate initialize mode
		initialize = false
		if SlamCfg.PrivEPK, err = epk.ReadFrom(fd); err != nil {
			return err
		}
		lm.GetLogger(`application`).Infoln(`successfully loaded private key from file`)
	}

	// test loaded credentials
	if SlamCfg.PubKey, err = SlamCfg.PrivEPK.Public(SlamCfg.Passphrase); err != nil {
		return err
	}
	lm.GetLogger(`application`).Infoln(`successfully unlocked public key from private key`)

	if initialize {
		if err = registerMachineEnrollment(); err != nil {
			return err
		}
	}

	return nil
}

func createKeypairFiles(path string) error {
	var err error
	var fd *os.File

	if SlamCfg.PrivEPK, SlamCfg.PubKey, err = epk.New(SlamCfg.Passphrase); err != nil {
		return err
	}

	if fd, err = os.Create(filepath.Join(path, `machinekey.epk`)); err != nil {
		return err
	}
	if err = SlamCfg.PrivEPK.Store(fd); err != nil {
		return err
	}
	fd.Close()

	if fd, err = os.Create(filepath.Join(path, `machinekey.pub`)); err != nil {
		return err
	}
	fmt.Fprintf(fd, "%s %s\n", `ed25519-epk-pub`, base64.StdEncoding.EncodeToString(SlamCfg.PubKey))
	fd.Close()
	return nil
}

func mask(pass []byte) []byte {
	pass = []byte(strings.TrimSpace(string(pass)))

	for i := range pass {
		// ROT13 ;)
		pass[i] = pass[i] | 0b00001101
		// mask used password vs disk stored password
		pass[i] = pass[i] ^ byte(0b01001101)

		pass[i] = bits.Reverse8(pass[i])
	}
	return pass
}

func createPassphraseFile(path string) error {
	d := make([]byte, 32)
	if _, err := io.ReadFull(cryptorand.Reader, d[:32]); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%x\n", d)
	f.Close()
	SlamCfg.Passphrase = string(mask([]byte(fmt.Sprintf("%x", d))))
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
