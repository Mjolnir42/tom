/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
)

const magicHashKey string = `engineroom.machine.tom`

func GetHash(pub ed25519.PublicKey) (string, error) {
	var err error
	var hfunc hash.Hash
	var dgst []byte

	if hfunc, err = blake2b.New(16, []byte(magicHashKey)); err != nil {
		return ``, err
	}

	hfunc.Write([]byte(pub))
	dgst = hfunc.Sum(nil)
	return fmt.Sprintf("%x", dgst), nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
