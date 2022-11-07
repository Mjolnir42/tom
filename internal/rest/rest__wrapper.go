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
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
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
		x.auth(
			func(w http.ResponseWriter, r *http.Request,
				ps httprouter.Params) {
				h(w, r, ps)
			},
		),
	)
}

// Deny is the request wrapper to straight up refuse processing
func (x *Rest) Deny(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}
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

// auth handles HTTP authentication on requests
func (x *Rest) auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		if !x.conf.Enforce {
			ps = append(ps, httprouter.Param{
				Key:   `AuthenticatedUser`,
				Value: `system~nobody`,
			})
			h(w, r, ps)
			return
		}

		auth := r.Header.Get(`Authorization`)
		switch {
		case strings.HasPrefix(auth, proto.AuthSchemeBasic+` `):
			x.basicAuth(
				func(w http.ResponseWriter, r *http.Request,
					ps httprouter.Params) {
					h(w, r, ps)
				},
			)
			return
		case strings.HasPrefix(auth, proto.AuthSchemeEPK+` `):
			x.epkAuth(
				func(w http.ResponseWriter, r *http.Request,
					ps httprouter.Params) {
					h(w, r, ps)
				},
			)
			return
		}

		w.Header().Set(`WWW-Authenticate`, `Basic realm=Restricted`)
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}
}

// basicAuth handles HTTP BasicAuth on requests
func (x *Rest) basicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		// Get credentials
		auth := r.Header.Get(`Authorization`)
		if strings.HasPrefix(auth, proto.AuthSchemeBasic+` `) {
			// check credentials
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(proto.AuthSchemeBasic+` `):])
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
		} else {
			// XXX no Authorization header configured
			ps = append(ps, httprouter.Param{
				Key:   `AuthenticatedUser`,
				Value: `system~nobody`,
			})
			h(w, r, ps)
			return
		}
	unauthorized:
		w.Header().Set(`WWW-Authenticate`, `Basic realm=Restricted`)
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}
}

// epkAuth handles HTTP signature on requests
func (x *Rest) epkAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var payload, nonceBytes, sigBytes []byte
		var err error
		var data []string
		var nonce, requestURI, idLib, userID, sig string
		var result msg.Result
		var request msg.Request

		// ...
		auth := r.Header.Get(`Authorization`)
		if !strings.HasPrefix(auth, proto.AuthSchemeEPK+` `) {
			goto unauthorized
		}

		if payload, err = base64.StdEncoding.DecodeString(auth[len(proto.AuthSchemeEPK+` `):]); err != nil {
			goto unauthorized
		}

		data = strings.Split(string(payload), `:`)
		nonce = data[0]
		requestURI = data[1]
		idLib = data[2]
		userID = data[3]
		sig = data[4]

		if nonceBytes, err = base64.StdEncoding.DecodeString(nonce); err != nil {
			goto unauthorized
		}
		if sigBytes, err = base64.StdEncoding.DecodeString(sig); err != nil {
			goto unauthorized
		}

		if requestURI != r.URL.Path {
			x.LM.GetLogger(`error`).Errorf("Mismatched request path in epkAuth: %s vs %s", requestURI, r.URL.Path)
			goto unauthorized
		}

		request = msg.New(
			r, ps,
			proto.CmdSupervisorAuthEPK,
			proto.EntitySupervisor,
			proto.ActionAuthenticateEPK,
		)
		request.Auth = msg.Super{
			Nonce:      nonceBytes,
			RequestURI: requestURI,
			IDLib:      idLib,
			UserID:     userID,
			Sig:        sigBytes,
		}
		x.HM.MustLookup(&request).Intake() <- request
		result = <-request.Reply
		if result.Err != nil {
			goto unauthorized
		}
		switch result.Code {
		case 200:
			ps = append(ps, httprouter.Param{
				Key:   `AuthenticatedUser`,
				Value: idLib + `~` + userID,
			})
			h(w, r, ps)
			return

		default:
			goto unauthorized
		}

	unauthorized:
		w.Header().Set(`WWW-Authenticate`, `Basic realm=Restricted`)
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
