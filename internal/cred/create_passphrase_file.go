/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/mjolnir42/tom/internal/config"
)

func createPassphraseFile(path string, cfg *config.AuthConfiguration) error {
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
	cfg.Passphrase = string(mask([]byte(fmt.Sprintf("%x", d))))
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
