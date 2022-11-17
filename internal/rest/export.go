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
	r.ApplyVerbosity()
	x.send(w, r, e)
}

func (x *Rest) IsAuthorized(q *msg.Request) bool {
	switch x.isAuthorized(q) {
	case true:
		x.LM.GetLogger(`audit`).Infof(
			"[RequestID] %s, [Section] %s, [Action] %s, [User] %s, [Code] %d",
			q.ID.String(), q.Section, q.Action, q.AuthUser, 200,
		)
		return true
	default:
		x.LM.GetLogger(`audit`).Errorf(
			"[RequestID] %s, [Section] %s, [Action] %s, [User] %s, [Code] %d",
			q.ID.String(), q.Section, q.Action, q.AuthUser, 403,
		)
		return false
	}
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

// Authorize is a simple authorization function
func Authorize(q *msg.Request) bool {
	// if authentication is not enforced, the authorization is always
	// false
	if !q.Enforcement {
		return false
	}

	// authentication is enforced, but the anonymous user appeared
	if q.AuthUser == `system~nobody` {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
