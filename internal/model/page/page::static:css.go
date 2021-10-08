/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package page // import "github.com/mjolnir42/tom/internal/model/page/"

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdPageStaticCSS)

	registry = append(registry, function{
		cmd:    proto.CmdPageStaticCSS,
		handle: staticCSS,
	})
}

func staticCSS(m *Model) httprouter.Handle {
	return m.x.Unauthenticated(m.StaticCSS)
}

// StaticCSS function
func (m *Model) StaticCSS(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	a, err := Asset(params.ByName(`asset`))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch strings.HasSuffix(params.ByName(`asset`), `.css`) {
	case true:
		w.Header().Add(`Content-Type`, `text/css`)
	default:
		w.Header().Add(`Content-Type`, `application/octet-stream`)
	}
	w.Header().Add(`X-Content-Type-Options`, `application/nosniff`)
	w.Write(a)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
