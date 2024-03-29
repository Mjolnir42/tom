/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"fmt"
	"net/http"

	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func nopExport(res *proto.Result, r *msg.Result) {
}

func (x *Rest) replyBadRequestDispatch(w *http.ResponseWriter, q *msg.Request, err error) {
	if err == nil {
		err = fmt.Errorf(http.StatusText(http.StatusBadRequest))
	}
	result := msg.FromRequest(q)
	result.BadRequest(err)
	x.send(w, &result, nopExport)
}

func (x *Rest) replyUnauthorizedDispatch(w *http.ResponseWriter, q *msg.Request) {
	result := msg.FromRequest(q)
	result.Unauthorized()
	x.send(w, &result, nopExport)
}

func (x *Rest) replyForbiddenDispatch(w *http.ResponseWriter, q *msg.Request) {
	result := msg.FromRequest(q)
	result.Forbidden()
	x.send(w, &result, nopExport)
}

func (x *Rest) replyServerError(w *http.ResponseWriter, q *msg.Request, err error) {
	if err == nil {
		err = fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	result := msg.FromRequest(q)
	result.ServerError()
	x.send(w, &result, nopExport)
}

func (x *Rest) hardServerError(w *http.ResponseWriter) {
	http.Error(*w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
