/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/tom/cmd/slamdd-go"

import (
	"os"
	"path/filepath"

	"github.com/droundy/goopt"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/sirupsen/logrus"
)

// global variables
var (
	// config file runtime configuration
	SlamCfg config.SlamConfiguration
	// lookup table of logfile handles for logrotate reopen
	lm *lhm.LogHandleMap
	// populated via Makefile
	slamVersion string
)

const (
	EX_OK    = 0
	EX_ERROR = 1
	EX_ABORT = 2
)

// startup initialization
func init() {
	lm = lhm.Init()
	// setup goopt information
	goopt.Version = slamVersion
	goopt.Suite = `slamdd`
	goopt.Summary = `slamDD`
	goopt.Author = `Jörg Pernfuß`
	goopt.Description = func() string {
		return "slam Data Daemon"
	}
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		err        error
		configFile string
	)
	//
	cliConfigPath := goopt.String([]string{`-c`, `--config`}, `./slam.conf`, `Configuration file location`)
	goopt.Parse(nil)

	// read configuration file
	lm.EarlyPrintf("Starting runtime config initialization, slamDD %s", slamVersion)
	if configFile, err = filepath.Abs(*cliConfigPath); err != nil {
		lm.EarlyFatal(err)
	}
	if configFile, err = filepath.EvalSymlinks(configFile); err != nil {
		lm.EarlyFatal(err)
	}
	lm.EarlyPrintf("Reading configuration file: %s", configFile)
	if err = SlamCfg.Parse(configFile, lm); err != nil {
		lm.EarlyFatal(err)
	}
	SlamCfg.Version = slamVersion
	lm.Setup(SlamCfg.LogPath)
	if err = lm.Open(`error`, logrus.ErrorLevel); err != nil {
		logrus.Fatal(err)
	}
	if err = lm.Open(`application`, logrus.InfoLevel); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}
	if err = lm.Open(`request`, logrus.InfoLevel); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}

	go lm.Reopen(``, func(e error) {
		logrus.Error(e)
		os.Exit(EX_ABORT)
	})

	lm.GetLogger(`application`).Infoln(`Loading credentials`)
	if err = loadCredentials(); err != nil {
		lm.GetLogger(`error`).Errorln(err)
		return EX_ERROR
	}

	return EX_OK
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
