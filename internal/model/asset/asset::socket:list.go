// +build socket

/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdSocketList)

	registry = append(registry, function{
		cmd:    proto.CmdSocketList,
		handle: socketList,
	})
}

func socketList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.SocketList)
}

func exportSocketList(result *proto.Result, r *msg.Result) {
	result.SocketHeader = &[]proto.SocketHeader{}
	*result.SocketHeader = append(*result.SocketHeader, r.SocketHeader...)
}

// SocketList function
func (m *Model) SocketList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// socket is uniquely identified. Process this as SocketShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.SocketShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionSocket
	request.Action = proto.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}
	if r.URL.Query().Get(`namespace`) != `` {
		request.Socket.Namespace = r.URL.Query().Get(`namespace`)
		if err := proto.ValidNamespace(request.Socket.Namespace); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportSocketList)
}

// list returns all servers
func (h *SocketReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		dictionaryName, socketName, author string
		creationTime                       time.Time
		namespace                          sql.NullString
		rows                               *sql.Rows
		err                                error
	)

	if q.Socket.Namespace != `` {
		namespace.String = q.Socket.Namespace
		namespace.Valid = true
	}

	if rows, err = h.stmtList.Query(
		namespace,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&dictionaryName,
			&socketName,
			&author,
			&creationTime,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		mr.SocketHeader = append(mr.SocketHeader, proto.SocketHeader{
			Namespace: dictionaryName,
			Name:      socketName,
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
