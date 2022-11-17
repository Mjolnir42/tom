/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
)

// Authorization ...
type Authorization struct {
	Timestamp   string     `json:"timestamp"`
	UserID      string     `json:"userID"`
	Fingerprint string     `json:"key.fingerprint"`
	Nonce       string     `json:"nonce"`
	Sig         *Signature `json:"signature,omitempty"`
	CSR         *DataCSR   `json:"csr,omitempty"`
}

// Serialize ...
func (a *Authorization) Serialize() []byte {
	if a.Nonce == `` {
		nonce := make([]byte, 6)
		io.ReadFull(cryptorand.Reader, nonce[:6])
		a.Nonce = base64.StdEncoding.EncodeToString(nonce)
	}

	data := make([]byte, 0)
	data = append(data, []byte(a.Timestamp)...)
	data = append(data, []byte(a.UserID)...)
	data = append(data, []byte(a.Fingerprint)...)
	nonceBytes, _ := base64.StdEncoding.DecodeString(a.Nonce)
	data = append(data, nonceBytes...)
	if a.CSR != nil {
		if a.CSR.Sig != nil {
			data = append(data, []byte(a.CSR.Sig.DataHash)...)
		}
	}
	return data
}

// Signature ...
type Signature struct {
	DataHash  string `json:"-"`
	Signature string `json:"signature"`
}

// DataCSR ...
type DataCSR struct {
	UserID       string     `json:"user-name"`
	Library      string     `json:"identity-library"`
	FQDN         string     `json:"fqdn"`
	PublicKey    string     `json:"public-key"`
	EnrolmentKey string     `json:"enrolment-key"`
	ValidFrom    string     `json:"valid-from"`
	ValidUntil   string     `json:"valid-until"`
	Sig          *Signature `json:"signature,omitempty"`
}

// Serialize ...
func (d *DataCSR) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(d.UserID)...)
	data = append(data, []byte(d.Library)...)
	data = append(data, []byte(d.FQDN)...)
	data = append(data, []byte(d.PublicKey)...)
	data = append(data, []byte(d.EnrolmentKey)...)
	data = append(data, []byte(d.ValidFrom)...)
	data = append(data, []byte(d.ValidUntil)...)
	return data
}

func (d *DataCSR) CalculateDataHash() error {
	var (
		err   error
		dgst  []byte
		hfunc hash.Hash
	)

	if hfunc, err = blake2b.New512(nil); err != nil {
		return err
	}
	hfunc.Write(d.Serialize())
	dgst = hfunc.Sum(nil)

	switch d.Sig {
	case nil:
		d.Sig = &Signature{
			DataHash: base64.StdEncoding.EncodeToString(dgst),
		}
	default:
		d.Sig.DataHash = base64.StdEncoding.EncodeToString(dgst)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
