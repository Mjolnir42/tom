/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/nahanni/go-ucl"
	"github.com/urfave/cli/v2"
)

type Config struct {
	API      string        `json:"api"`
	LogDir   string        `json:"logdir"`
	ProcJSON string        `json:"json.output.processor"`
	Run      RunTimeConfig `json:"-"`
}

type RunTimeConfig struct {
	API      *url.URL `json:"-"`
	PathLogs string   `json:"-"`
}

func (c *Config) populateFromFile(fname string) error {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	// UCL parses into map[string]interface{}
	fileBytes := bytes.NewBuffer([]byte(file))
	parser := ucl.NewParser(fileBytes)
	uclData, err := parser.Ucl()
	if err != nil {
		return err
	}

	// take detour via JSON to load UCL into struct
	uclJSON, err := json.Marshal(uclData)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(uclJSON), &c)

	return nil
}

func configSetup(c *cli.Context) (*Config, error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var confPath string
	switch {
	case c.IsSet(`config`) && path.IsAbs(c.String(`config`)):
		confPath = path.Clean(c.String(`config`))
	case c.IsSet(`config`):
		confPath = path.Clean(path.Join(home, ".tom", c.String(`config`)))
	default:
		confPath = path.Clean(path.Join(home, ".tom", "tom.conf"))
	}

	cfg := &Config{Run: RunTimeConfig{}}
	if err = cfg.populateFromFile(confPath); err != nil {
		return nil, err
	}
	cfg.Run.API, err = url.Parse(cfg.API)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
