/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package config // import "github.com/mjolnir42/tom/internal/config"

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/lhm"
	"github.com/nahanni/go-ucl"
	"golang.org/x/crypto/ed25519"
)

type SlamConfiguration struct {
	Daemon     []Daemon                 `json:"daemon"`
	LogLevel   string                   `json:"log.level"`
	LogPath    string                   `json:"log.path"`
	CredPath   string                   `json:"credential.path"`
	IPFIX      SettingsIPFIX            `json:"settings.ipfix"`
	Version    string                   `json:"-"`
	Passphrase string                   `json:"-"`
	PubKey     ed25519.PublicKey        `json:"-"`
	PrivEPK    *epk.EncryptedPrivateKey `json:"-"`
}

type SettingsIPFIX struct {
	Enabled      bool   `json:"enabled,string"`
	ServerProto  string `json:"listen.protocol"`
	ListenADDR   string `json:"listen.address"`
	Forwarding   bool   `json:"forwarding.enabled,string"`
	ForwardADDR  string `json:"forwarding.address"`
	ForwardProto string `json:"forwarding.protocol"`
	ServerName   string `json:"tls.servername"`
	CAFile       string `json:"ca.file"`
	CertFile     string `json:"certificate.file"`
	CertKeyFile  string `json:"certificate.keyfile"`
	Processing   bool   `json:"processing.enabled,string"`
	ProcessType  string `json:"processing.type"`
}

func (c *SlamConfiguration) Parse(fname string, lh *lhm.LogHandleMap) error {
	//
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	lh.EarlyPrintf("Loading configuration from %s", fname)

	// UCL parses into map[string]interface{}
	fileBytes := bytes.NewBuffer([]byte(file))
	parser := ucl.NewParser(fileBytes)
	uclData, err := parser.Ucl()
	if err != nil {
		lh.EarlyFatal("UCL error: ", err)
	}

	// take detour via JSON to load UCL into struct
	uclJSON, err := json.Marshal(uclData)
	if err != nil {
		lh.EarlyFatal(err)
	}
	json.Unmarshal([]byte(uclJSON), &c)

	//
	switch c.LogLevel {
	case `debug`, `info`, `warn`, `error`, `fatal`, `panic`:
	default:
		lh.EarlyFatal(`Invalid log.level specified: `, c.LogLevel, `. Valid levels are: `,
			`debug, info (default), warn, error, fatal, panic`)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
