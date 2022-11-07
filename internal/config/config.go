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

	"github.com/mjolnir42/lhm"
	"github.com/nahanni/go-ucl"
)

type Configuration struct {
	Database DbConfig   `json:"database"`
	Daemon   []Daemon   `json:"daemon"`
	Auth     AuthConfig `json:"authentication"`
	LogLevel string     `json:"log.level"`
	LogPath  string     `json:"log.path"`
	Version  string     `json:"-"`
	QueueLen int        `json:"handler.queue.length,string"`
	Enforce  bool       `json:"enforcement,string"`
}

type DbConfig struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Name    string `json:"database"`
	Port    string `json:"port"`
	Pass    string `json:"password"`
	Timeout string `json:"timeout"`
	TLSMode string `json:"tlsmode"`
}

type Daemon struct {
	URL    *url.URL `json:"-"`
	Listen string   `json:"listen"`
	Port   string   `json:"port"`
	TLS    bool     `json:"tls,string"`
	Cert   string   `json:"cert.file"`
	Key    string   `json:"key.file"`
}

func (c *Configuration) Parse(fname string, lh *lhm.LogHandleMap) error {
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

	if c.QueueLen <= 0 {
		lh.EarlyPrintf("Adjusting QueueLen from %d to %d (default)", c.QueueLen, 8)
		c.QueueLen = 8
	}

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
