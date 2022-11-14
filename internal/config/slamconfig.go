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
	Auth     *AuthConfiguration `json:"authentication"`
	IPFIX    SettingsIPFIX      `json:"ipfix"`
	LogLevel string             `json:"log.level"`
	LogPath  string             `json:"log.path"`
	API      string             `json:"api"`
	CAFile   string             `json:"api.ca.file"`
	Version  string             `json:"-"`
	Run      RunTimeConfig      `json:"-"`
}

type AuthConfiguration struct {
	Passphrase string                   `json:"-"`
	PubKey     ed25519.PublicKey        `json:"-"`
	PrivEPK    *epk.EncryptedPrivateKey `json:"-"`
	CredPath   string                   `json:"credential.path"`
}

type SettingsIPFIX struct {
	Enabled     bool       `json:"enabled,string"`
	Forwarding  bool       `json:"forwarding.enabled,string"`
	Processing  bool       `json:"processing.enabled,string"`
	ProcessType string     `json:"processing.type"`
	TemplFile   string     `json:"template.file"`
	Servers     []IPDaemon `json:"server"`
	Clients     []IPClient `json:"client"`
	Filters     IPFilter   `json:"filter"`
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

type IPClient struct {
	Enabled      bool   `json:"enabled,string"`
	ForwardADDR  string `json:"forwarding.address"`
	ForwardProto string `json:"forwarding.protocol"`
	CAFile       string `json:"ca.file"`
	Unfiltered   bool   `json:"unfiltered.copy,string"`
	Format       string `json:"json.format"` // vflow,flowdata
}

type IPFilter struct {
	Rules  []string `json:"rules"`
	Parsed []Rule   `json:"-"`
}

type Rule struct {
	MatchField         string
	FieldType          string
	MatchValueString   []string
	MatchValueUint8    []uint8
	MatchValueUint16   []uint16
	Action             string
	ReplaceValueString string
	ReplaceValueUint8  uint8
	ReplaceValueUint16 uint16
	InverseMatch       bool
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

	c.IPFIX.ProcessType = strings.TrimSpace(c.IPFIX.ProcessType)
	c.IPFIX.TemplFile = strings.TrimSpace(c.IPFIX.TemplFile)

	for i := range c.IPFIX.Servers {
		c.IPFIX.Servers[i].CAFile = strings.TrimSpace(c.IPFIX.Servers[i].CAFile)
		c.IPFIX.Servers[i].CAFile = filepath.Clean(c.IPFIX.Servers[i].CAFile)
		c.IPFIX.Servers[i].CertFile = strings.TrimSpace(c.IPFIX.Servers[i].CertFile)
		c.IPFIX.Servers[i].CertFile = filepath.Clean(c.IPFIX.Servers[i].CertFile)
		c.IPFIX.Servers[i].CertKeyFile = strings.TrimSpace(c.IPFIX.Servers[i].CertKeyFile)
		c.IPFIX.Servers[i].CertKeyFile = filepath.Clean(c.IPFIX.Servers[i].CertKeyFile)
		c.IPFIX.Servers[i].ListenADDR = strings.TrimSpace(c.IPFIX.Servers[i].ListenADDR)
		c.IPFIX.Servers[i].ServerName = strings.TrimSpace(c.IPFIX.Servers[i].ServerName)
		c.IPFIX.Servers[i].ServerProto = strings.TrimSpace(c.IPFIX.Servers[i].ServerProto)
	}
	for i := range c.IPFIX.Clients {
		c.IPFIX.Clients[i].CAFile = strings.TrimSpace(c.IPFIX.Clients[i].CAFile)
		c.IPFIX.Clients[i].CAFile = filepath.Clean(c.IPFIX.Clients[i].CAFile)
		c.IPFIX.Clients[i].ForwardADDR = strings.TrimSpace(c.IPFIX.Clients[i].ForwardADDR)
		c.IPFIX.Clients[i].ForwardProto = strings.TrimSpace(c.IPFIX.Clients[i].ForwardProto)
	}

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
