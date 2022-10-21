/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"encoding/base64"
	"hash"

	"golang.org/x/crypto/blake2b"
)

// Request is the request wrapper of Tom's public API
type Request struct {
	Verbose       bool           `json:"verbose,omitempty,string"`
	Container     *Container     `json:"container,omitempty"`
	Library       *Library       `json:"library,omitempty"`
	Namespace     *Namespace     `json:"namespace,omitempty"`
	Orchestration *Orchestration `json:"orchestration,omitempty"`
	Runtime       *Runtime       `json:"runtime,omitempty"`
	Server        *Server        `json:"server,omitempty"`
	Socket        *Socket        `json:"socket,omitempty"`
	Team          *Team          `json:"team,omitempty"`
	User          *User          `json:"user,omitempty"`
	Auth          Authorization  `json:"authorization"`
}

// Serialize ...
func (r *Request) Serialize() []byte {
	data := make([]byte, 0)
	if r.Container != nil {
		data = append(data, r.Container.Serialize()...)
	}
	if r.Library != nil {
		data = append(data, r.Library.Serialize()...)
	}
	if r.Namespace != nil {
		data = append(data, r.Namespace.Serialize()...)
	}
	if r.Orchestration != nil {
		data = append(data, r.Orchestration.Serialize()...)
	}
	if r.Runtime != nil {
		data = append(data, r.Runtime.Serialize()...)
	}
	if r.Server != nil {
		data = append(data, r.Server.Serialize()...)
	}
	if r.Socket != nil {
		data = append(data, r.Socket.Serialize()...)
	}
	if r.Team != nil {
		data = append(data, r.Team.Serialize()...)
	}
	if r.User != nil {
		data = append(data, r.User.Serialize()...)
	}
	data = append(data, r.Auth.Serialize()...)
	return data
}

// CalculateDataHash ...
func (r *Request) CalculateDataHash() error {
	var (
		err   error
		dgst  []byte
		hfunc hash.Hash
	)

	if hfunc, err = blake2b.New512(nil); err != nil {
		return err
	}
	hfunc.Write(r.Serialize())
	dgst = hfunc.Sum(nil)

	r.Auth.Sig = &Signature{
		DataHash: base64.StdEncoding.EncodeToString(dgst),
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
