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
	"net/url"
	"path/filepath"
	"strings"

	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/lhm"
	"github.com/nahanni/go-ucl"
	"golang.org/x/crypto/ed25519"
)

type SlamConfiguration struct {
	Daemon   []Daemon           `json:"daemon"`
	LogLevel string             `json:"log.level"`
	LogPath  string             `json:"log.path"`
	IPFIX    SettingsIPFIX      `json:"ipfix"`
	IPFIXSrv []IPDaemon         `json:"ipfix.server"`
	API      string             `json:"api"`
	CAFile   string             `json:"api.ca.file"`
	Auth     *AuthConfiguration `json:"authentication"`
	Version  string             `json:"-"`
	Run      RunTimeConfig      `json:"-"`
}

type AuthConfiguration struct {
	Passphrase string                   `json:"-"`
	PubKey     ed25519.PublicKey        `json:"-"`
	PrivEPK    *epk.EncryptedPrivateKey `json:"-"`
	CredPath   string                   `json:"credential.path"`
}

type IPDaemon struct {
	Enabled     bool   `json:"enabled,string"`
	ServerProto string `json:"listen.protocol"`
	ListenADDR  string `json:"listen.address"`
	ServerName  string `json:"tls.servername"`
	CAFile      string `json:"ca.file"`
	CertFile    string `json:"certificate.file"`
	CertKeyFile string `json:"certificate.keyfile"`
}

type SettingsIPFIX struct {
	Enabled      bool   `json:"enabled,string"`
	Forwarding   bool   `json:"forwarding.enabled,string"`
	ForwardADDR  string `json:"forwarding.address"`
	ForwardProto string `json:"forwarding.protocol"`
	CAFile       string `json:"ca.file"`
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

	if c.Run.API, err = url.Parse(c.API); err != nil {
		return err
	}

	if c.CAFile != `` {
		c.Run.PathCA = filepath.Clean(c.CAFile)
	}

	c.IPFIX.CAFile = strings.TrimSpace(c.IPFIX.CAFile)
	c.IPFIX.ForwardADDR = strings.TrimSpace(c.IPFIX.ForwardADDR)
	c.IPFIX.ForwardProto = strings.TrimSpace(c.IPFIX.ForwardProto)
	c.IPFIX.ProcessType = strings.TrimSpace(c.IPFIX.ProcessType)

	for i := range c.IPFIXSrv {
		c.IPFIXSrv[i].CAFile = strings.TrimSpace(c.IPFIXSrv[i].CAFile)
		c.IPFIXSrv[i].CertFile = strings.TrimSpace(c.IPFIXSrv[i].CertFile)
		c.IPFIXSrv[i].CertKeyFile = strings.TrimSpace(c.IPFIXSrv[i].CertKeyFile)
		c.IPFIXSrv[i].ListenADDR = strings.TrimSpace(c.IPFIXSrv[i].ListenADDR)
		c.IPFIXSrv[i].ServerName = strings.TrimSpace(c.IPFIXSrv[i].ServerName)
		c.IPFIXSrv[i].ServerProto = strings.TrimSpace(c.IPFIXSrv[i].ServerProto)
	}

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
