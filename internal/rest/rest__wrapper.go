/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/tom/internal/rest/"

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

// Unauthenticated is a wrapper for unauthenticated or implicitly
// authenticated requests
func (x *Rest) Unauthenticated(h httprouter.Handle) httprouter.Handle {
	return x.enrich(
		func(w http.ResponseWriter, r *http.Request,
			ps httprouter.Params) {
			h(w, r, ps)
		},
	)
}

// Authenticated is the standard request wrapper
func (x *Rest) Authenticated(h httprouter.Handle) httprouter.Handle {
	return x.Unauthenticated(
		x.basicAuth(
			func(w http.ResponseWriter, r *http.Request,
				ps httprouter.Params) {
				h(w, r, ps)
			},
		),
	)
}

// enrich is a wrapper that adds metadata information to the request
func (x *Rest) enrich(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		// generate and record the requestID
		requestID := uuid.NewV4()
		ps = append(ps, httprouter.Param{
			Key:   `RequestID`,
			Value: requestID.String(),
		})

		// record the request URI
		ps = append(ps, httprouter.Param{
			Key:   `RequestURI`,
			Value: r.RequestURI,
		})

		h(w, r, ps)
	}
}

// basicAuth handles HTTP BasicAuth on requests
func (x *Rest) basicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "

		// Get credentials
		auth := r.Header.Get(`Authorization`)
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// check credentials
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 {
					// TODO test credentials
					// TODO if !OK goto unauthorized
					ps = append(ps, httprouter.Param{
						Key:   `AuthenticatedUser`,
						Value: string(pair[0]),
					})
					// record the used token for supervisor:token/invalidate
					ps = append(ps, httprouter.Param{
						Key:   `AuthenticatedToken`,
						Value: string(pair[1]),
					})

					// Delegate request to handle
					h(w, r, ps)
					return
				} else {
					goto unauthorized
				}
			}
		}
	unauthorized:
		w.Header().Set(`WWW-Authenticate`, `Basic realm=Restricted`)
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
