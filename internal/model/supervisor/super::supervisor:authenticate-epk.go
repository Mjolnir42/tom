/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package supervisor // import "github.com/mjolnir42/tom/internal/model/supervisor/"

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/cred"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
	"golang.org/x/crypto/ed25519"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdSupervisorAuthEPK)

	registry = append(registry, function{
		cmd:    proto.CmdSupervisorAuthEPK,
		handle: supervisorAuthEPK,
	})
}

// REST calls into the supervisor are invalid
func supervisorAuthEPK(m *Model) httprouter.Handle {
	return m.x.Deny(m.SupervisorAuth)
}

// SupervisorAuth function
func (m *Model) SupervisorAuth(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// wrapped in Rest.Deny() this function should never be reached
	request := msg.New(
		r, params,
		proto.CmdSupervisorAuthEPK,
		proto.EntitySupervisor,
		proto.ActionAuthenticateEPK,
	)
	m.x.ReplyForbidden(&w, &request)
}

// authenticateEPK ...
func (h *CoreHandler) authenticateEPK(q *msg.Request, mr *msg.Result) {
	var (
		data                                      []string
		nonce, unixT, reqURI, fp, idLib, uID, sig string
		fingerprint, publicKey                    string
		tx                                        *sql.Tx
		err                                       error
		unix                                      int64
		pubKey                                    []byte
	)

	// request.Auth.Token split -> base64/nonce
	//                          -> string/unixTime
	//                          -> string/requestURI
	//                          -> string/fingerprint
	//                          -> string/IDLib
	//                          -> string/userID
	//                          -> base64/signature
	if data = strings.Split(q.Auth.Token, `:`); len(data) < 7 {
		mr.ExpectationFailed(fmt.Errorf(
			"Incomplete TOM-epk authentication token of length %d",
			len(data)),
		)
		return
	}
	nonce = data[0]
	unixT = data[1]
	// the requestURI might contain : characters
	reqURI = strings.Join(data[2:len(data)-4], `:`)
	fp = data[len(data)-4]
	idLib = data[len(data)-3]
	uID = data[len(data)-2]
	sig = data[len(data)-1]

	// check if requestURI paths match
	if reqURI != q.Auth.RequestURI {
		mr.ServerError(fmt.Errorf(
			"Mismatching requestURI paths: %s vs %s",
			reqURI, q.Auth.RequestURI),
		)
		return
	}

	if unix, err = strconv.ParseInt(unixT, 10, 64); err != nil {
		mr.ServerError(err)
		return
	}
	txTime := time.Now().UTC()
	switch {
	case time.Unix(unix, 0).UTC().After(txTime):
		fallthrough
	case time.Unix(unix, 0).UTC().Before(txTime.Add(-30 * time.Second)):
		mr.ExpectationFailed(fmt.Errorf(
			"TOM-epk signature outside validity: %s",
			txTime.Sub(time.Unix(unix, 0).UTC()).String(),
		))
	default:
	}

	// start transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}
	if _, err = tx.Exec(stmt.ReadOnlyTransaction); err != nil {
		mr.ServerError(err)
		return
	}
	defer tx.Rollback()

	txKey := tx.Stmt(h.stmtKey)

	if err = txKey.QueryRow(
		uID,
		idLib,
		txTime,
	).Scan(
		&publicKey,
		&fingerprint,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	if fp != fingerprint {
		mr.ServerError(fmt.Errorf(
			"Mismatching signing key fingerprints: %s vs %s",
			fp, fingerprint,
		))
		return
	}

	pack := msg.Super{
		//Nonce: Base64DecodeString(nonce)
		FP:         fingerprint,
		RequestURI: q.Auth.RequestURI,
		IDLib:      idLib,
		UserID:     uID,
	}
	pack.Time = make([]byte, 8)
	binary.BigEndian.PutUint64(pack.Time, uint64(unix))

	pack.Nonce, err = base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		mr.ServerError(err)
		return
	}
	pack.Signature, err = base64.StdEncoding.DecodeString(sig)
	if err != nil {
		mr.ServerError(err)
		return
	}
	pubKey, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		mr.ServerError(err)
		return
	}
	pack.Public = ed25519.PublicKey(pubKey)

	if !cred.VerifyEpkAuthToken(pack) {
		mr.Unauthorized()
		return
	}

	mr.Auth.IDLib = idLib
	mr.Auth.UserID = uID
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
