/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"fmt"
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
	Err        error

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

func (r *Result) Clear(err ...error) {
	if len(err) > 0 {
		r.Err = err[len(err)-1]
	}

	switch r.Section {
	case SectionServer:
		r.Server = []proto.Server{}
	}
}

func (r *Result) OK() {
	r.Code = http.StatusOK
}

func (r *Result) Forbidden() {
	r.Code = http.StatusForbidden
	r.Clear()
}

func (r *Result) ServerError(err ...error) {
	r.Code = http.StatusInternalServerError
	r.Clear()
}

func (r *Result) NotImplemented(err ...error) {
	r.Code = http.StatusNotImplemented
	r.Clear()
}

func (r *Result) UnknownRequest(rq *Request) {
	r.NotImplemented(fmt.Errorf("Unknown requested action:"+
		" %s/%s", rq.Section, rq.Action))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
