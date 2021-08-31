/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespaceAdd)

	registry = append(registry, function{
		cmd:    proto.CmdNamespaceAdd,
		method: (proto.Commands[proto.CmdNamespaceAdd]).Method,
		path:   (proto.Commands[proto.CmdNamespaceAdd]).Path,
		handle: namespaceAdd,
	})
}

func namespaceAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespaceAdd)
}

func exportNamespaceAdd(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespaceAdd function
func (m *Model) NamespaceAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = *req.Namespace

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportNamespaceAdd)
}

// add creates a new namespace
func (h *NamespaceWriteHandler) add(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
		tx  *sql.Tx
	)

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// create named namespace
	if res, err = tx.Stmt(h.stmtAdd).Exec(
		q.Namespace.Property[`dict_name`].Value,
		q.UserIDLib,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if !mr.CheckRowsAffected(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	// set configured namespace properties
	for _, property := range []string{
		`dict_type`,
		`dict_lookup`,
		`dict_uri`,
		`dict_ntt_list`,
	} {
		if _, ok := q.Namespace.Property[property]; ok {
			if res, err = tx.Stmt(h.stmtConfig).Exec(
				q.Namespace.Property[`dict_name`].Value,
				property,
				q.Namespace.Property[property].Value,
				q.UserIDLib,
				q.AuthUser,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
			if !mr.CheckRowsAffected(res.RowsAffected()) {
				tx.Rollback()
				return
			}
		}
	}

	// create configured namespace attributes
	for _, attribute := range q.Namespace.Attributes {
		if attribute.Unique {
			res, err = tx.Stmt(h.stmtAttUnqAdd).Exec(
				q.Namespace.Property[`dict_name`].Value,
				attribute.Key,
				q.UserIDLib,
				q.AuthUser,
			)
		} else {
			res, err = tx.Stmt(h.stmtAttStdAdd).Exec(
				q.Namespace.Property[`dict_name`].Value,
				attribute.Key,
				q.UserIDLib,
				q.AuthUser,
			)
		}
		if err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		if !mr.CheckRowsAffected(res.RowsAffected()) {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
