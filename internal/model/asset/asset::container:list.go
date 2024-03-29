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
	proto.AssertCommandIsDefined(proto.CmdContainerList)

	registry = append(registry, function{
		cmd:    proto.CmdContainerList,
		handle: containerList,
	})
}

func containerList(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ContainerList)
}

func exportContainerList(result *proto.Result, r *msg.Result) {
	result.ContainerHeader = &[]proto.ContainerHeader{}
	*result.ContainerHeader = append(*result.ContainerHeader, r.ContainerHeader...)
}

// ContainerList function
func (m *Model) ContainerList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// if both ?name and ?namespace are set as query paramaters, the
	// container is uniquely identified. Process this as ContainerShow request
	if r.URL.Query().Get(`name`) != `` && r.URL.Query().Get(`namespace`) != `` {
		m.ContainerShow(w, r, params)
		return
	}

	request := msg.New(
		r, params,
		proto.CmdContainerList,
		msg.SectionContainer,
		proto.ActionList,
	)

	if r.URL.Query().Get(`namespace`) != `` {
		request.Container.Namespace = r.URL.Query().Get(`namespace`)
		if err := proto.ValidNamespace(request.Container.Namespace); err != nil {
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
	m.x.Send(&w, &result, exportContainerList)
}

// list returns all containers
func (h *ContainerReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		cnrID, nsName, key, value, author string
		creationTime                      time.Time
		namespace                         sql.NullString
		rows                              *sql.Rows
		container                         proto.ContainerHeader
		err                               error
		ok                                bool
	)

	if q.Container.Namespace != `` {
		namespace.String = q.Container.Namespace
		namespace.Valid = true
	}

	list := make(map[string]proto.ContainerHeader)
	if rows, err = h.stmtList.Query(
		namespace,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&cnrID,
			&nsName,
			&key,
			&value,
			&author,
			&creationTime,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		if container, ok = list[cnrID]; !ok {
			container = proto.ContainerHeader{}
		}
		container.Namespace = nsName
		switch {
		case key == `type`:
			container.Type = value
		case key == `name`:
			container.Name = value
			container.CreatedBy = author
			container.CreatedAt = creationTime.Format(msg.RFC3339Milli)
		}
		list[cnrID] = container
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	for _, container := range list {
		mr.ContainerHeader = append(mr.ContainerHeader, container)
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
