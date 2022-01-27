/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

type ExportFunc func(pub *proto.Result, internal *msg.Result)

type Rest struct {
	isAuthorized func(*msg.Request) bool
	conf         *config.Configuration
	HM           *handler.Map
	LM           *lhm.LogHandleMap
	idx          int
}

func New(
	authorizationFunction func(*msg.Request) bool,
	index int,
	hm *handler.Map,
	lm *lhm.LogHandleMap,
	cfg *config.Configuration,
) *Rest {
	x := Rest{}
	x.isAuthorized = authorizationFunction
	x.idx = index
	x.HM = hm
	x.LM = lm
	x.conf = cfg
	return &x
}

func (x *Rest) Run(rt *httprouter.Router) {
	// TODO switch to new abortable interface
	switch x.conf.Daemon[x.idx].URL.Scheme {
	case `http`:
		x.LM.GetLogger(`error`).Fatal(http.ListenAndServe(
			x.conf.Daemon[x.idx].URL.Host,
			rt,
		))
	case `https`:
		x.LM.GetLogger(`error`).Fatal(http.ListenAndServeTLS(
			x.conf.Daemon[x.idx].URL.Host,
			x.conf.Daemon[x.idx].Cert,
			x.conf.Daemon[x.idx].Key,
			rt,
		))
	default:
		x.LM.GetLogger(`error`).Fatalf(
			"Unsupported URL scheme: %s",
			x.conf.Daemon[x.idx].URL.Scheme,
		)
	}
}

func PanicCatcher(w http.ResponseWriter, lm *lhm.LogHandleMap) {
	if r := recover(); r != nil {
		lm.GetLogger(`error`).Errorf("PANIC! %s, TRACE: %s", r, debug.Stack())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
}

// decodeJSONBody unmarshals the JSON request body from r into s
func DecodeJSONBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
		if err == nil {
			if c.Container != nil {
				if c.Container.Link == nil {
					c.Container.Link = []string{}
				}
				if c.Container.Property == nil {
					c.Container.Property = map[string]proto.PropertyDetail{}
				}
			}
			if c.Orchestration != nil {
				if c.Orchestration.Parent == nil {
					c.Orchestration.Parent = []string{}
				}
				if c.Orchestration.Link == nil {
					c.Orchestration.Link = []string{}
				}
				if c.Orchestration.Children == nil {
					c.Orchestration.Children = []string{}
				}
				if c.Orchestration.Property == nil {
					c.Orchestration.Property = map[string]proto.PropertyDetail{}
				}
			}
			if c.Runtime != nil {
				if c.Runtime.Link == nil {
					c.Runtime.Link = []string{}
				}
				if c.Runtime.Children == nil {
					c.Runtime.Children = []string{}
				}
				if c.Runtime.Property == nil {
					c.Runtime.Property = map[string]proto.PropertyDetail{}
				}
			}
			if c.Server != nil {
				if c.Server.Link == nil {
					c.Server.Link = []string{}
				}
				if c.Server.Children == nil {
					c.Server.Children = []string{}
				}
				if c.Server.Property == nil {
					c.Server.Property = map[string]proto.PropertyDetail{}
				}
			}
		}
	default:
		rt := reflect.TypeOf(s)
		err = fmt.Errorf("DecodeJSON: unhandled request of type: %s", rt)
	}
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
