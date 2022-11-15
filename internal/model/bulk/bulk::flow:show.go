/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package bulk // import "github.com/mjolnir42/tom/internal/model/bulk/"

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
	proto.AssertCommandIsDefined(proto.CmdFlowShow)

	registry = append(registry, function{
		cmd:    proto.CmdFlowShow,
		handle: flowShow,
	})
}

func flowShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.FlowShow)
}

func exportFlowShow(result *proto.Result, r *msg.Result) {
	result.Flow = &[]proto.Flow{}
	*result.Flow = append(*result.Flow, r.Flow...)
}

// FlowShow function
func (m *Model) FlowShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdFlowShow,
		msg.SectionFlow,
		proto.ActionShow,
	)
	request.Flow.TomID = params.ByName(`tomID`)
	request.Flow.Namespace = r.URL.Query().Get(`namespace`)
	request.Flow.Name = r.URL.Query().Get(`name`)

	if err := request.Flow.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
		// error is ErrEmptyTomID, check for query parameter
		// supplied values
		switch {
		case request.Flow.Namespace == ``:
			fallthrough
		case request.Flow.Name == ``:
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
	m.x.Send(&w, &result, exportFlowShow)
}

// show returns full details for a specific flow
func (h *FlowReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		flowID, dictionaryID string
		tx                   *sql.Tx
		rows                 *sql.Rows
		err                  error
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

	flow := *(proto.NewFlow())
	flow.Namespace = q.Flow.Namespace
	flow.Name = q.Flow.Name
	name := proto.PropertyDetail{
		Attribute: `name`,
		Value:     q.Flow.Name,
	}

	var since, until, createdAt, namedAt time.Time
	if err = txShow.QueryRow(
		q.Flow.Namespace,
		q.Flow.Name,
		txTime,
	).Scan(
		&flowID,
		&dictionaryID,
		&createdAt,
		&flow.CreatedBy,
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

	flow.CreatedAt = createdAt.Format(msg.RFC3339Milli)
	name.CreatedAt = namedAt.Format(msg.RFC3339Milli)
	name.ValidSince = since.Format(msg.RFC3339Milli)
	name.ValidUntil = until.Format(msg.RFC3339Milli)
	name.Namespace = q.Flow.Namespace
	flow.Property = make(map[string]proto.PropertyDetail)
	flow.Property[q.Flow.Namespace+`::`+flow.Name+`::name`] = name

	// fetch flow properties
	if rows, err = txProp.Query(
		dictionaryID,
		flowID,
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
		prop.ValidSince = since.Format(msg.RFC3339Milli)
		prop.ValidUntil = until.Format(msg.RFC3339Milli)
		prop.CreatedAt = at.Format(msg.RFC3339Milli)
		prop.Namespace = q.Flow.Namespace

		// set specialty fields
		switch prop.Attribute {
		case `name`:
			if flow.Name != prop.Value {
				rows.Close()
				mr.ExpectationFailed(
					fmt.Errorf(`Encountered confused resultset`),
				)
				return
			}
		case `type`:
			flow.Type = prop.Value
		}
		switch prop.Attribute {
		case `name`:
			// name attribute has already been added
		default:
			flow.Property[prop.Namespace+`::`+flow.Name+`::`+prop.Attribute] = prop
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
	mr.Flow = append(mr.Flow, flow)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
