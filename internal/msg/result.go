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

	Container           []proto.Container
	ContainerHeader     []proto.ContainerHeader
	Library             []proto.Library
	Namespace           []proto.Namespace
	NamespaceHeader     []proto.NamespaceHeader
	Orchestration       []proto.Orchestration
	OrchestrationHeader []proto.OrchestrationHeader
	Runtime             []proto.Runtime
	RuntimeHeader       []proto.RuntimeHeader
	Server              []proto.Server
	ServerHeader        []proto.ServerHeader
	Team                []proto.Team
	User                []proto.User
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

func (r *Result) CheckRowsAffected(i int64, err error) bool {
	if err != nil {
		r.ServerError(err)
		return false
	}
	switch i {
	case 0:
		r.OK()
		return true
	case 1:
		r.OK()
		return true
	default:
		r.ServerError(fmt.Errorf("Too many rows affected: %d", i))
		return false
	}
}

func (r *Result) AssertOneRowAffected(i int64, err error) bool {
	if err != nil {
		r.ServerError(err)
		return false
	}
	switch i {
	case 1:
		r.OK()
		return true
	default:
		r.ServerError(fmt.Errorf("Assertion failed: %d rows affected instead of 1", i))
		return false
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
