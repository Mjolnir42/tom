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
	Command    string
	Section    string
	Action     string
	Code       uint16
	Err        error
	Verbose    bool

	Auth Super

	Container           []proto.Container
	ContainerHeader     []proto.ContainerHeader
	Flow                []proto.Flow
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
		Command:    rq.Command,
		RequestURI: rq.RequestURI,
		Section:    rq.Section,
		Action:     rq.Action,
		Code:       http.StatusNotImplemented,
		Verbose:    rq.Verbose,
	}
}

func (r *Result) Clear(err ...error) {
	if len(err) > 0 {
		r.Err = err[len(err)-1]
	}
	if r.Err == nil {
		r.Err = fmt.Errorf(`Unspecified error condition`)
	}

	switch r.Section {
	case SectionNamespace:
		r.Namespace = []proto.Namespace{}
		r.NamespaceHeader = []proto.NamespaceHeader{}
	case SectionOrchestration:
		r.Orchestration = []proto.Orchestration{}
		r.OrchestrationHeader = []proto.OrchestrationHeader{}
	case SectionRuntime:
		r.Runtime = []proto.Runtime{}
		r.RuntimeHeader = []proto.RuntimeHeader{}
	case SectionServer:
		r.Server = []proto.Server{}
		r.ServerHeader = []proto.ServerHeader{}
	case SectionContainer:
		r.Container = []proto.Container{}
		r.ContainerHeader = []proto.ContainerHeader{}
	case SectionFlow:
		r.Flow = []proto.Flow{}
	case proto.EntitySupervisor:
		r.Auth = Super{}
	}
}

func (r *Result) OK() {
	r.Code = http.StatusOK
}

func (r *Result) BadRequest(err ...error) {
	r.Code = http.StatusBadRequest
	r.Clear(err...)
}

func (r *Result) Unauthorized(err ...error) {
	r.Code = http.StatusUnauthorized
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

func (r *Result) ExportError() error {
	return r.Err
}

func (r *Result) ApplyVerbosity() {
	if r.Verbose {
		return
	}
	for i := range r.Container {
		obj := r.Container[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		for k := range obj.Property {
			prop := obj.Property[k]
			prop.CreatedAt = ``
			prop.CreatedBy = ``
			prop.ValidSince = ``
			prop.ValidUntil = ``
			obj.Property[k] = prop
		}
		r.Container[i] = obj
	}
	for i := range r.ContainerHeader {
		obj := r.ContainerHeader[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.ContainerHeader[i] = obj
	}
	for i := range r.Flow {
		obj := r.Flow[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		for k := range obj.Property {
			prop := obj.Property[k]
			prop.CreatedAt = ``
			prop.CreatedBy = ``
			prop.ValidSince = ``
			prop.ValidUntil = ``
			obj.Property[k] = prop
		}
		r.Flow[i] = obj
	}
	for i := range r.Library {
		obj := r.Library[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.Library[i] = obj
	}
	for i := range r.Namespace {
		obj := r.Namespace[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		for k := range obj.Property {
			prop := obj.Property[k]
			prop.CreatedAt = ``
			prop.CreatedBy = ``
			prop.ValidSince = ``
			prop.ValidUntil = ``
			obj.Property[k] = prop
		}
		r.Namespace[i] = obj
	}
	for i := range r.NamespaceHeader {
		obj := r.NamespaceHeader[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.NamespaceHeader[i] = obj
	}
	for i := range r.Orchestration {
		obj := r.Orchestration[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		for k := range obj.Property {
			prop := obj.Property[k]
			prop.CreatedAt = ``
			prop.CreatedBy = ``
			prop.ValidSince = ``
			prop.ValidUntil = ``
			obj.Property[k] = prop
		}
		r.Orchestration[i] = obj
	}
	for i := range r.OrchestrationHeader {
		obj := r.OrchestrationHeader[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.OrchestrationHeader[i] = obj
	}
	for i := range r.Runtime {
		obj := r.Runtime[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		for k := range obj.Property {
			prop := obj.Property[k]
			prop.CreatedAt = ``
			prop.CreatedBy = ``
			prop.ValidSince = ``
			prop.ValidUntil = ``
			obj.Property[k] = prop
		}
		r.Runtime[i] = obj
	}
	for i := range r.RuntimeHeader {
		obj := r.RuntimeHeader[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.RuntimeHeader[i] = obj
	}
	for i := range r.Server {
		obj := r.Server[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		for k := range obj.Property {
			prop := obj.Property[k]
			prop.CreatedAt = ``
			prop.CreatedBy = ``
			prop.ValidSince = ``
			prop.ValidUntil = ``
			obj.Property[k] = prop
		}
		r.Server[i] = obj
	}
	for i := range r.ServerHeader {
		obj := r.ServerHeader[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.ServerHeader[i] = obj
	}
	for i := range r.Team {
		obj := r.Team[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.Team[i] = obj
	}
	for i := range r.User {
		obj := r.User[i]
		obj.CreatedAt = ``
		obj.CreatedBy = ``
		r.User[i] = obj
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
