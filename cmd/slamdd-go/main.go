/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/tom/cmd/slamdd-go"

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/droundy/goopt"
	"github.com/go-resty/resty/v2"
	"github.com/mjolnir42/epk"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/cli/adm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/cred"
	"github.com/mjolnir42/tom/internal/ipfix"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

// global variables
var (
	// config file runtime configuration
	SlamCfg config.SlamConfiguration
	// lookup table of logfile handles for logrotate reopen
	lm *lhm.LogHandleMap
	// populated via Makefile
	slamVersion string
	//
	client *resty.Client
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
		session    tls.ClientSessionCache
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

	if SlamCfg.Enabled {
		if SlamCfg.Auth == nil {
			SlamCfg.Auth = &config.AuthConfiguration{}
		}
		SlamCfg.Auth.PrivEPK = &epk.EncryptedPrivateKey{}
		SlamCfg.Auth.PubKey = ed25519.PublicKey{}

		// setup REST client
		client = resty.New().
			SetDisableWarn(true).
			SetHeader(`User-Agent`, fmt.Sprintf("%s %s", goopt.Summary, goopt.Version)).
			SetHostURL(SlamCfg.Run.API.String())
		if SlamCfg.Run.API.Scheme == `https` {
			session = tls.NewLRUClientSessionCache(64)

			client = client.SetTLSClientConfig(&tls.Config{
				ServerName:         strings.SplitN(SlamCfg.Run.API.Host, `:`, 2)[0],
				ClientSessionCache: session,
				MinVersion:         tls.VersionTLS12,
			}).SetRootCertificate(SlamCfg.Run.PathCA)
		}

		adm.ConfigureClient(client)

		var initialize bool
		lm.GetLogger(`application`).Infoln(`Loading credentials`)
		if initialize, err = cred.LoadCredentials(
			SlamCfg.Auth,
			lm,
		); err != nil {
			lm.GetLogger(`error`).Errorln(err)
			return EX_ERROR
		} else if initialize {
			lm.GetLogger(`application`).Infoln(`registering newly initialized credentials with TOM service`)
			adm.RegisterMachineEnrollment(SlamCfg.Auth, nil)
		}

	}
	var ipfEx chan interface{}
	if SlamCfg.IPFIX.Enabled {
		if ipfEx, err = ipfix.New(SlamCfg, lm); err != nil {
			lm.GetLogger(`error`).Errorln(err)
			return EX_ERROR
		}
	}

	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)

	lm.GetLogger(`application`).Println(`Waiting for signals`)
	select {
	case <-ipfEx:
	case <-cancel:
		// give ipfix subsystem a chance to shut down gracefully
		select {
		case <-time.After(7 * time.Second):
		case <-ipfEx:
		}
	}

	return EX_OK
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
