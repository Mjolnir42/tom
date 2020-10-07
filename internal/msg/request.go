/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/pkg/proto"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	ID         uuid.UUID
	Section    string
	Action     string
	RemoteAddr string
	AuthUser   string
	RequestURI string
	Reply      chan Result `json:"-"`

	Server proto.Server
}

// New returns a Request
func New(r *http.Request, params httprouter.Params) Request {
	returnChannel := make(chan Result, 1)
	return Request{
		ID:         requestID(params),
		RequestURI: requestURI(params),
		RemoteAddr: remoteAddr(r),
		AuthUser:   authUser(params),
		Reply:      returnChannel,
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
