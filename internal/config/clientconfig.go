/*-
 * Copyright (c) 2021, Jörg Pernfuß
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

	"github.com/nahanni/go-ucl"
)

type ClientConfig struct {
	API      string             `json:"api"`
	LogDir   string             `json:"logdir"`
	ProcJSON string             `json:"json.output.processor"`
	CAFile   string             `json:"ca.file"`
	Auth     *AuthConfiguration `json:"authentication"`
	Run      RunTimeConfig      `json:"-"`
}

type RunTimeConfig struct {
	API      *url.URL `json:"-"`
	PathLogs string   `json:"-"`
	PathCA   string   `json:"-"`
}

func (c *ClientConfig) PopulateFromFile(fname string) error {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
