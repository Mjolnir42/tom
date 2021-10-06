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

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdServerList)

	registry = append(registry, function{
		cmd:    proto.CmdServerList,
		handle: serverList,
	})
}

func serverList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ServerList)
}

func exportServerList(result *proto.Result, r *msg.Result) {
	result.ServerHeader = &[]proto.ServerHeader{}
	*result.ServerHeader = append(*result.ServerHeader, r.ServerHeader...)
}

// ServerList function
func (m *Model) ServerList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// server is uniquely identified. Process this as ServerShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.ServerShow(w, r, params)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = proto.ActionList

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportServerList)
}

// list returns all servers
func (h *ServerReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		id, namespace, key, value string
		rows                      *sql.Rows
		err                       error
		server                    proto.Server
		ok                        bool
	)

	list := make(map[string]proto.Server)
	if rows, err = h.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&namespace,
			&key,
			&value,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		if server, ok = list[id]; !ok {
			server = proto.Server{}
		}
		server.ID = id
		server.Namespace = namespace
		switch {
		case key == `type`:
			server.Type = value
		case key == `name`:
			server.Name = value
		}
		list[id] = server
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	for _, server := range list {
		mr.Server = append(mr.Server, server)
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix