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
	proto.AssertCommandIsDefined(proto.CmdMachEnrol)

	registry = append(registry, function{
		cmd:    proto.CmdMachEnrol,
		handle: machineEnrol,
	})
}

func machineEnrol(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.MachineEnrol)
}

func exportMachineEnrol(result *proto.Result, r *msg.Result) {
	result.User = &[]proto.User{}
	*result.User = append(*result.User, r.User...)
}

// MachineEnrol function
func (m *Model) MachineEnrol(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdMachEnrol,
		msg.SectionMachine,
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
			"Invalid credential type for machine registrations: %s",
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

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportMachineEnrol)
}

// enrolment is the handler for registering machine accounts
func (h *UserWriteHandler) enrolment(q *msg.Request, mr *msg.Result) {
	var (
		err                        error
		isMachine, isSelfEnrolment bool
		enrolmentKey               sql.NullString
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

	if !isMachine {
		mr.ExpectationFailed(
			fmt.Errorf(`Only registered machine libraries support machine enrolment`),
		)
		return
	}

	switch {
	case isSelfEnrolment:
	case !isSelfEnrolment && enrolmentKey.Valid:
	case !isSelfEnrolment && !enrolmentKey.Valid:
		mr.ExpectationFailed(
			fmt.Errorf(`Machine library without self-enrolment requires enrolment key`),
		)
		return
	default:
	}

	// create registration
	if err = tx.Stmt(h.stmtEnrolment).QueryRow(
		q.User.LibraryID,
		q.User.FirstName,
		q.User.LastName,
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
			fmt.Errorf(`Machine registration should contain credential of type public-key`),
		)
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
