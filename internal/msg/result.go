/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"net/http"

	"github.com/mjolnir42/tom/pkg/proto"
	uuid "github.com/satori/go.uuid"
)

type Result struct {
	ID         uuid.UUID
	RequestURI string
	Section    string
	Action     string
	Code       uint16

	Server []proto.Server
}

func FromRequest(rq *Request) Result {
	return Result{
		ID:         rq.ID,
		RequestURI: rq.RequestURI,
		Section:    rq.Section,
		Action:     rq.Action,
		Code:       http.StatusNotImplemented,
	}
}

func (r *Result) OK() {
	r.Code = http.StatusOK
}

func (r *Result) Forbidden() {
	r.Code = http.StatusForbidden
}

func (r *Result) ServerError() {
	r.Code = http.StatusInternalServerError
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
