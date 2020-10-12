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

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

type Rest struct {
	isAuthorized func(*msg.Request) bool
	conf         *config.Configuration
	hm           *handler.Map
	lm           *lhm.LogHandleMap
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
	x.hm = hm
	x.lm = lm
	x.conf = cfg
	return &x
}

func (x *Rest) Run() {
	router := x.setupRouter()

	// TODO switch to new abortable interface
	x.lm.GetLogger(`error`).Fatal(http.ListenAndServeTLS(
		x.conf.Daemon[x.idx].URL.Host,
		x.conf.Daemon[x.idx].Cert,
		x.conf.Daemon[x.idx].Key,
		router,
	))
}

func panicCatcher(w http.ResponseWriter, lm *lhm.LogHandleMap) {
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
func decodeJSONBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		err = fmt.Errorf("DecodeJSON: unhandled request of type: %s", rt)
	}
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
