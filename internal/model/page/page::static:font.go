/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package page // import "github.com/mjolnir42/tom/internal/model/page/"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdPageStaticFont)

	registry = append(registry, function{
		cmd:    proto.CmdPageStaticFont,
		handle: staticFont,
	})
}

func staticFont(m *Model) httprouter.Handle {
	return m.x.Unauthenticated(m.StaticFont)
}

// StaticFont function
func (m *Model) StaticFont(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	a, err := Asset(params.ByName(`asset`))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add(`Content-Type`, `application/font-woff2`)
	w.Write(a)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
