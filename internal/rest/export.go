/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"net/http"

	"github.com/mjolnir42/tom/internal/msg"
)

func (x *Rest) Send(w *http.ResponseWriter, r *msg.Result, e ExportFunc) {
	x.send(w, r, e)
}

const (
	MethodDELETE = `DELETE`
	MethodGET    = `GET`
	MethodHEAD   = `HEAD`
	MethodPATCH  = `PATCH`
	MethodPOST   = `POST`
	MethodPUT    = `PUT`
)

func (x *Rest) IsAuthorized(q *msg.Request) bool {
	return x.isAuthorized(q)
}

func (x *Rest) ReplyBadRequest(w *http.ResponseWriter, q *msg.Request, err error) {
	x.replyBadRequestDispatch(w, q, err)
}

func (x *Rest) ReplyForbidden(w *http.ResponseWriter, q *msg.Request) {
	x.replyForbiddenDispatch(w, q)
}

func (x *Rest) ReplyServerError(w *http.ResponseWriter, q *msg.Request, err error) {
	x.replyServerError(w, q, err)
}

func (x *Rest) HardServerError(w *http.ResponseWriter) {
	x.hardServerError(w)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
