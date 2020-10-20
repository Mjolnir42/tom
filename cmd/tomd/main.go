/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/tom/"

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/droundy/goopt"
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/core"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/model/asset"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/sirupsen/logrus"
)

// global variables
var (
	// config file runtime configuration
	TomCfg config.Configuration
	// lookup table of logfile handles for logrotate reopen
	lm *lhm.LogHandleMap
	//
	tomVersion string
)

// startup initialization
func init() {
	lm = lhm.Init()
	// setup goopt information
	goopt.Version = tomVersion
	goopt.Suite = `tom`
	goopt.Summary = `Tom`
	goopt.Author = `Jörg Pernfuß`
	goopt.Description = func() string {
		return "Tom is the guy you ask if you need to know details about stuff."
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
	// parse command line options
	cliConfigPath := goopt.String([]string{`-c`, `--config`}, `./tom.conf`, `Configuration file location`)
	goopt.Parse(nil)

	// read configuration file
	lm.EarlyPrintf("Starting runtime config initialization, TOM v%s", tomVersion)
	if configFile, err = filepath.Abs(*cliConfigPath); err != nil {
		lm.EarlyFatal(err)
	}
	if configFile, err = filepath.EvalSymlinks(configFile); err != nil {
		lm.EarlyFatal(err)
	}
	lm.EarlyPrintf("Reading configuration file: %s", configFile)
	if err = TomCfg.Parse(configFile, lm); err != nil {
		lm.EarlyFatal(err)
	}
	TomCfg.Version = tomVersion
	lm.Setup(TomCfg.LogPath)
	if err = lm.Open(`error`, logrus.ErrorLevel); err != nil {
		logrus.Fatal(err)
	}
	if err = lm.Open(`application`, logrus.InfoLevel); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}

	go lm.Reopen(``, func(e error) {
		logrus.Error(e)
		os.Exit(2)
	})

	// create handler map
	hm := handler.NewMap()

	// start main database connection pool
	conn := connectToDatabase(lm)
	go pingDatabase(lm, conn)

	// start core application
	core := core.New(hm, lm, conn, &TomCfg)
	core.Start()

	for i := range TomCfg.Daemon {
		dm := TomCfg.Daemon[i]
		dm.URL = &url.URL{}
		dm.URL.Host = fmt.Sprintf("%s:%s", dm.Listen, dm.Port)
		if dm.TLS {
			dm.URL.Scheme = `https`
		} else {
			dm.URL.Scheme = `http`
		}
		api := rest.New(func(q *msg.Request) bool { return true }, i, hm, lm, &TomCfg)
		router := httprouter.New()

		// create datamodels
		assetmodel := asset.New(api)
		router = assetmodel.RouteRegisterServer(router)
		router = assetmodel.RouteRegisterRuntime(router)

		go api.Run(router)
	}

	// signal handler for shutdown
	sigChanShutdown := make(chan os.Signal, 1)
	signal.Notify(sigChanShutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-sigChanShutdown
	return 0
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
