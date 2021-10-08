/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package page // import "github.com/mjolnir42/tom/internal/model/page/"

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdPageApplication)

	registry = append(registry, function{
		cmd:    proto.CmdPageApplication,
		handle: application,
	})
}

func application(m *Model) httprouter.Handle {
	return m.x.Unauthenticated(m.Application)
}

// Application function
func (m *Model) Application(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	tmpl, err := Asset(`tmpl/base.html`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	t, err := template.New(`base`).Parse(string(tmpl))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	t.Execute(w, nil)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
