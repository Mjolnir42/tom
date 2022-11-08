/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"strconv"
	"time"

	"github.com/mjolnir42/tom/internal/msg"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
)

// CalcEpkAuthToken ...
func CalcEpkAuthToken(data msg.Super) (string, error) {
	var (
		err                         error
		hfunc                       hash.Hash
		tstp                        time.Time
		sig, dgst, nonce, timeBytes []byte
		token, fp                   string
		pubkey                      ed25519.PublicKey
	)

	// generate nonce
	nonce = make([]byte, 6)
	io.ReadFull(cryptorand.Reader, nonce[:6])

	// generate time
	tstp = time.Now().UTC()
	timeBytes = make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(tstp.Unix()))

	// generate pubkey
	if pubkey, err = data.PK.Public(data.Phrase); err != nil {
		return ``, err
	}

	// generate key fingerprint
	if fp, err = GetHash(pubkey); err != nil {
		return ``, err
	}

	// generate digest
	if hfunc, err = blake2b.New(16, nil); err != nil {
		return ``, err
	}
	hfunc.Write(nonce)
	hfunc.Write(timeBytes)
	hfunc.Write([]byte(fp))
	hfunc.Write([]byte(data.RequestURI))
	hfunc.Write([]byte(data.IDLib))
	hfunc.Write([]byte(data.UserID))
	dgst = hfunc.Sum(nil)

	// sign digest
	if sig, err = data.PK.Sign(data.Phrase, dgst); err != nil {
		return ``, err
	}

	// generate token data
	token = fmt.Sprintf(
		"%s:%s:%s:%s:%s:%s:%s",
		base64.StdEncoding.EncodeToString(nonce),
		strconv.FormatInt(tstp.Unix(), 10),
		data.RequestURI,
		fp,
		data.IDLib,
		data.UserID,
		base64.StdEncoding.EncodeToString(sig),
	)
	// return base64 encoded token
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

// VerifyEpkAuthToken ...
func VerifyEpkAuthToken(data msg.Super) bool {
	var (
		err   error
		hfunc hash.Hash
		dgst  []byte
	)

	// generate digest
	if hfunc, err = blake2b.New(16, nil); err != nil {
		return false
	}
	hfunc.Write(data.Nonce)
	hfunc.Write(data.Time)
	hfunc.Write([]byte(data.FP))
	hfunc.Write([]byte(data.RequestURI))
	hfunc.Write([]byte(data.IDLib))
	hfunc.Write([]byte(data.UserID))
	dgst = hfunc.Sum(nil)

	return ed25519.Verify(data.Public, dgst, data.Signature)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
