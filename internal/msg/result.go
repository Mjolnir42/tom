/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
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

	Namespace       []proto.Namespace
	NamespaceHeader []proto.NamespaceHeader
	Orchestration   []proto.Orchestration
	Runtime         []proto.Runtime
	Server          []proto.Server
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
	case SectionNamespace:
		r.Namespace = []proto.Namespace{}
		r.NamespaceHeader = []proto.NamespaceHeader{}
	case SectionOrchestration:
		r.Orchestration = []proto.Orchestration{}
	case SectionRuntime:
		r.Runtime = []proto.Runtime{}
	case SectionServer:
		r.Server = []proto.Server{}
	}
}

func (r *Result) OK() {
	r.Code = http.StatusOK
}

func (r *Result) BadRequest(err ...error) {
	r.Code = http.StatusBadRequest
	r.Clear(err...)
}

func (r *Result) NotFound(err ...error) {
	r.Code = http.StatusNotFound
	r.Clear(err...)
}

func (r *Result) Forbidden(err ...error) {
	r.Code = http.StatusForbidden
	r.Clear(err...)
}

func (r *Result) ExpectationFailed(err ...error) {
	r.Code = http.StatusExpectationFailed
	r.Clear(err...)
}

func (r *Result) ServerError(err ...error) {
	r.Code = http.StatusInternalServerError
	r.Clear(err...)
}

func (r *Result) NotImplemented(err ...error) {
	r.Code = http.StatusNotImplemented
	r.Clear(err...)
}

func (r *Result) UnknownRequest(rq *Request) {
	r.NotImplemented(fmt.Errorf("Unknown requested action:"+
		" %s/%s", rq.Section, rq.Action))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
