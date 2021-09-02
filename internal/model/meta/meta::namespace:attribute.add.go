/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespaceAttrAdd)

	registry = append(registry, function{
		cmd:    proto.CmdNamespaceAttrAdd,
		handle: namespaceAttrAdd,
	})
}

func namespaceAttrAdd(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespaceAttrAdd)
}

func exportNamespaceAttrAdd(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespaceAttrAdd function
func (m *Model) NamespaceAttrAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionAttrAdd

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = *req.Namespace
	request.Namespace.TomID = params.ByName(`tomID`)
	if err := request.Namespace.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	for _, attribute := range request.Namespace.Attributes {
		if err := proto.OnlyUnreserved(attribute.Key); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}

		if strings.HasPrefix(attribute.Key, `dict_`) {
			switch attribute.Key {
			case `dict_lookup`:
			case `dict_uri`:
			case `dict_ntt_list`:
			default:
				m.x.ReplyBadRequest(&w, &request, fmt.Errorf(
					"Invalid namespace self-attribute: %s",
					attribute.Key,
				))
				return
			}
		}
	}

	if err := proto.OnlyUnreserved(request.Namespace.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportNamespaceAttrAdd)
}

// attributeAdd ...
func (h *NamespaceWriteHandler) attributeAdd(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
		tx  *sql.Tx
	)

	//
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	//
	for _, attribute := range q.Namespace.Attributes {
		if strings.HasPrefix(attribute.Key, `dict_`) {
			switch attribute.Key {
			case `dict_lookup`:
				fallthrough
			case `dict_uri`:
				fallthrough
			case `dict_ntt_list`:
				if attribute.Unique {
					mr.BadRequest(fmt.Errorf(
						"Invalid unique namespace self-attribute: %s",
						attribute.Key,
					))
					tx.Rollback()
					return
				}
			default:
				mr.BadRequest(fmt.Errorf(
					"Invalid namespace self-attribute: %s",
					attribute.Key,
				))
				tx.Rollback()
				return
			}
		}

		if attribute.Unique {
			res, err = tx.Stmt(h.stmtAttUnqAdd).Exec(
				q.Namespace.Name,
				attribute.Key,
				q.UserIDLib,
				q.AuthUser,
			)
		} else {
			res, err = tx.Stmt(h.stmtAttStdAdd).Exec(
				q.Namespace.Name,
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
