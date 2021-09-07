/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"fmt"
	"net/http"
	//	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespacePropUpdate)

	registry = append(registry, function{
		cmd:    proto.CmdNamespacePropUpdate,
		handle: namespacePropertyUpdate,
	})
}

func namespacePropertyUpdate(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespacePropertyUpdate)
}

func exportNamespacePropertyUpdate(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespacePropertyUpdate function
func (m *Model) NamespacePropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = proto.ActionPropUpdate

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

	if err := proto.ValidNamespace(request.Namespace.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	// check property structure is setup
	if request.Namespace.Property == nil {
		m.x.ReplyBadRequest(&w, &request, fmt.Errorf(
			"Invalid %s request without properties",
			proto.ActionPropUpdate,
		))
		return
	}

	// check property names are all valid
	for prop := range request.Namespace.Property {
		if err := proto.OnlyUnreserved(prop); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportNamespacePropertyUpdate)
}

// propertyUpdate updates the specified properties of a namespace.
func (h *NamespaceWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	var (
		err    error
		ok     bool
		res    sql.Result
		rows   *sql.Rows
		tx     *sql.Tx
		txTime time.Time
	)
	// setup a consistent transaction time timestamp that is used for all
	// record updates during the transaction to make them concurrent
	txTime = time.Now().UTC()

	// tx.Begin()
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover all attributes and record them with their type
	attrMap := map[string]string{}
	if rows, err = tx.Stmt(h.stmtAttDiscover).Query(
		q.Namespace.Name,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	for rows.Next() {
		var attribute, typ string
		if err = rows.Scan(
			&attribute,
			&typ,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		attrMap[attribute] = typ
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// for all properties specified in the request, check that the attribute
	// exists and create missing attributes
	for key := range q.Namespace.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = tx.Stmt(h.stmtAttStdAdd).Exec(
				q.Namespace.Name,
				key,
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

	// for all properties specified in the request, update the value.
	// this transparently creates missing entries.
	for key := range q.Namespace.Property {
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Namespace.Property[key],
		); !ok {
			tx.Rollback()
			return
		}
	}
	// tx.Commit()
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
