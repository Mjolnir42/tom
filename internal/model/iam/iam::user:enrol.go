/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
	"golang.org/x/crypto/blake2b"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdUserEnrol)

	registry = append(registry, function{
		cmd:    proto.CmdUserEnrol,
		handle: userEnrol,
	})
}

func userEnrol(m *Model) httprouter.Handle {
	return m.x.Unauthenticated(m.UserEnrol)
}

func exportUserEnrol(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// UserEnrol function
func (m *Model) UserEnrol(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdUserEnrol,
		msg.SectionUser,
		proto.ActionEnrolment,
	)

	req := proto.NewUserRequest()
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.User = *req.User
	request.User.TomID = params.ByName(`tomID`)

	if err := request.User.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if err := proto.OnlyUnreserved(request.User.UserName); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	switch request.User.Credential.Category {
	case proto.CredentialPubKey:
	default:
		m.x.ReplyBadRequest(&w, &request, fmt.Errorf(
			"Invalid credential type for user enrolment: %s",
			request.User.Credential.Category,
		))
		return
	}

	switch request.User.Credential.Category {
	case proto.CredentialPubKey:
		if ok, err := req.Verify(); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		} else if !ok {
			m.x.LM.GetLogger(`application`).Infoln(
				request.ID, `Request signature verification failed`,
				req.Auth.Sig.DataHash,
				req.Auth.Sig.Signature,
			)
			m.x.ReplyForbidden(&w, &request)
			return
		}
	default:
		if !m.x.IsAuthorized(&request) {
			m.x.ReplyForbidden(&w, &request)
			return
		}
	}
	request.Auth = msg.Super{
		Token: req.Auth.CSR.EnrolmentKey,
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportUserEnrol)
}

// userEnrolment is the handler for registering user accounts
func (h *UserWriteHandler) userEnrolment(q *msg.Request, mr *msg.Result) {
	var (
		err                        error
		isMachine, isSelfEnrolment bool
		rootEnrolment              bool
		enrolmentKey, name         sql.NullString
		tx                         *sql.Tx
		res                        sql.Result
		hfunc                      hash.Hash
		key, dgst                  []byte
	)

	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}
	defer tx.Rollback()

	if _, err = tx.Exec(stmt.DeferredTransaction); err != nil {
		mr.ServerError(err)
		return
	}

	if err = tx.Stmt(h.stmtDetect).QueryRow(
		q.User.LibraryName,
	).Scan(
		&q.User.LibraryID,
		&isSelfEnrolment,
		&isMachine,
		&enrolmentKey,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	if isMachine {
		mr.ExpectationFailed(
			fmt.Errorf(`Attempted user enrolment to machine library`),
		)
		return
	}

	switch {
	case isSelfEnrolment:
		mr.ExpectationFailed(
			fmt.Errorf(`User identity library does not support self-enrolment`),
		)
		return
	case !isSelfEnrolment && !enrolmentKey.Valid:
		// there is a special case, where user enrolment without an
		// enrolment key is accepted:
		// - LibraryName: system
		// - UserName:    root
		// - Authentication Enforcement: false
		switch {
		case q.User.LibraryName != `system`:
			fallthrough
		case q.User.UserName != `root`:
			fallthrough
		case q.Enforcement:
			mr.ExpectationFailed(
				fmt.Errorf(`User library enrolment requires enrolment key`),
			)
			return
		default:
			rootEnrolment = true
			q.User.ID = `00000000-0000-0000-0000-000000000000`
		}
	case !isSelfEnrolment && enrolmentKey.Valid:
	default:
		mr.ServerError()
		return
	}

	if enrolmentKey.String != q.Auth.Token {
		mr.BadRequest()
		return
	}

	if !rootEnrolment {
		// create registration
		if err = tx.Stmt(h.stmtEnrolment).QueryRow(
			q.User.LibraryID,
			name,
			name,
			q.User.UserName,
			`ffffffff-ffff-ffff-ffff-ffffffffffff`,
		).Scan(
			&q.User.ID,
		); err == sql.ErrNoRows {
			mr.ServerError(err)
			return
		} else if err != nil {
			mr.ServerError(err)
			return
		}

		if res, err = tx.Stmt(h.stmtUpdateUID).Exec(
			q.User.ID,
			q.User.LibraryID,
		); err != nil {
			mr.ServerError(err)
			return
		}
		if !mr.AssertOneRowAffected(res.RowsAffected()) {
			return
		}
	}

	switch q.User.Credential.Category {
	case proto.CredentialPubKey:
		// calculate key fingerprint
		if hfunc, err = blake2b.New(16, []byte(`engineroom.machine.tom`)); err != nil {
			mr.ServerError(err)
			return
		}
		if key, err = base64.StdEncoding.DecodeString(q.User.Credential.Value); err != nil {
			mr.ServerError(err)
			return
		}
		hfunc.Write(key)
		dgst = hfunc.Sum(nil)

		// add key credential
		if res, err = tx.Stmt(h.stmtAddKey).Exec(
			q.User.ID,
			q.User.Credential.Value,
			fmt.Sprintf("%x", dgst),
			q.User.UserName,
			q.User.LibraryName,
			q.User.ID,
		); err != nil {
			mr.ServerError(err)
			return
		}
		if !mr.AssertOneRowAffected(res.RowsAffected()) {
			return
		}

	default:
		mr.ExpectationFailed(
			fmt.Errorf(`User enrolment should contain credential of type public-key`),
		)
		return
	}

	// ensure user record is active
	if res, err = tx.Stmt(h.stmtActivate).Exec(
		q.User.UserName,
		q.User.LibraryName,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	q.User.IsActive = true
	q.User.Credential = nil
	mr.User = append(mr.User, q.User)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
