/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdContainerShow)

	registry = append(registry, function{
		cmd:    proto.CmdContainerShow,
		handle: containerShow,
	})
}

func containerShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ContainerShow)
}

func exportContainerShow(result *proto.Result, r *msg.Result) {
	result.Container = &[]proto.Container{}
	*result.Container = append(*result.Container, r.Container...)
}

// ContainerShow function
func (m *Model) ContainerShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionContainer
	request.Action = proto.ActionShow
	request.Container = proto.Container{
		TomID:     params.ByName(`tomID`),
		Namespace: r.URL.Query().Get(`namespace`),
		Name:      r.URL.Query().Get(`name`),
	}

	if err := request.Container.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
		// error is ErrEmptyTomID, check for query parameter
		// supplied values
		switch {
		case request.Container.Namespace == ``:
			fallthrough
		case request.Container.Name == ``:
			m.x.ReplyBadRequest(&w, &request, nil)
			return
		}
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportContainerShow)
}

// show returns full details for a specific container
func (h *ContainerReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		rteID, dictionaryID string
		tx                  *sql.Tx
		rows                *sql.Rows
		err                 error
	)

	// start transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}
	if _, err = tx.Exec(stmt.ReadOnlyTransaction); err != nil {
		mr.ServerError(err)
		return
	}
	defer tx.Rollback()

	txTime := time.Now().UTC()
	txShow := tx.Stmt(h.stmtShow)
	txProp := tx.Stmt(h.stmtProp)

	rte := proto.Container{
		Namespace: q.Container.Namespace,
		Name:      q.Container.Name,
	}
	name := proto.PropertyDetail{
		Attribute: `name`,
		Value:     q.Container.Name,
	}

	var since, until, createdAt, namedAt time.Time
	if err = txShow.QueryRow(
		q.Container.Namespace,
		q.Container.Name,
		txTime,
	).Scan(
		&rteID,
		&dictionaryID,
		&createdAt,
		&rte.CreatedBy,
		&since,
		&until,
		&namedAt,
		&name.CreatedBy,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	rte.CreatedAt = createdAt.Format(msg.RFC3339Milli)
	name.CreatedAt = namedAt.Format(msg.RFC3339Milli)
	name.ValidSince = since.Format(msg.RFC3339Milli)
	name.ValidUntil = until.Format(msg.RFC3339Milli)
	rte.Property = make(map[string]proto.PropertyDetail)
	rte.Property[`name`] = name

	if rows, err = txProp.Query(
		dictionaryID,
		rteID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		prop := proto.PropertyDetail{}
		var since, until, at time.Time

		if err = rows.Scan(
			&prop.Attribute,
			&prop.Value,
			&since,
			&until,
			&at,
			&prop.CreatedBy,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		if prop.Attribute == `name` {
			// name has already been set
			continue
		}
		prop.ValidSince = since.Format(msg.RFC3339Milli)
		prop.ValidUntil = until.Format(msg.RFC3339Milli)
		prop.CreatedAt = at.Format(msg.RFC3339Milli)
		prop.Namespace = q.Container.Namespace

		// set specialty fields
		switch prop.Attribute {
		case `name`:
			if rte.Name != prop.Value {
				rows.Close()
				mr.ExpectationFailed(
					fmt.Errorf(`Encountered confused resultset`),
				)
				return
			}
		case `type`:
			rte.Type = prop.Value
		}
		switch prop.Attribute {
		case `name`:
			// name attribute has already been added
		default:
			rte.Property[prop.Attribute] = prop
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Container = append(mr.Container, rte)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
