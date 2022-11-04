/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
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

	req := proto.Request{}
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

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
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

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.User = append(mr.User, q.User)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
