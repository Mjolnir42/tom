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
	proto.AssertCommandIsDefined(proto.CmdPageStaticImage)

	registry = append(registry, function{
		cmd:    proto.CmdPageStaticImage,
		handle: staticImage,
	})
}

func staticImage(m *Model) httprouter.Handle {
	return m.x.Unauthenticated(m.StaticImage)
}

// StaticImage function
func (m *Model) StaticImage(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	a, err := Asset(params.ByName(`asset`))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add(`Content-Type`, `image/png`)
	w.Write(a)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
