/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	registry = append(registry, function{
		cmd:    proto.CmdNamespaceList,
		method: proto.MethodGET,
		path:   `/namespace/`,
		handle: namespaceList,
	})
}

func namespaceList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespaceList)
}

func exportNamespaceList(result *proto.Result, r *msg.Result) {
	result.NamespaceHeader = &[]proto.NamespaceHeader{}
	*result.NamespaceHeader = append(*result.NamespaceHeader, r.NamespaceHeader...)
}

// NamespaceList function
func (m *Model) NamespaceList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if ?name is set as query paramaters, the namespace is uniquely
	// identified. Process this as NamespaceShow request
	if r.URL.Query().Get(`name`) != `` {
		m.NamespaceShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportNamespaceList)
}

// list returns all namespaces
func (h *NamespaceReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		dictionaryName, author string
		creationTime           time.Time
		rows                   *sql.Rows
		err                    error
	)

	if rows, err = h.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&dictionaryName,
			&creationTime,
			&author,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		mr.NamespaceHeader = append(mr.NamespaceHeader, proto.NamespaceHeader{
			Name:      dictionaryName,
			CreatedAt: creationTime.Format(msg.RFC3339Milli),
			CreatedBy: author,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
