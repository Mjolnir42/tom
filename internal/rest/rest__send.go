/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"encoding/json"
	"net/http"

	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func (x *Rest) send(w *http.ResponseWriter, r *msg.Result) {
	var bjson []byte
	var err error
	var result proto.Result

	if r.Code >= http.StatusBadRequest {
		goto dispatchERROR
	}

	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)

	if bjson, err = json.Marshal(&result); err != nil {
		x.hardServerError(w)
		return
	}
	goto sendJSON

dispatchERROR:
	x.writeError(w, r.Code)
	return

sendJSON:
	x.writeReplyJSON(w, &bjson)
	return
}

// writeReplyJSON writes out b as the reply with content-type
// set to application/json
func (x *Rest) writeReplyJSON(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set(`Content-Type`, `application/json`)
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

func (x *Rest) writeError(w *http.ResponseWriter, code uint16) {
	(*w).WriteHeader(int(code))
	(*w).Write(nil)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
