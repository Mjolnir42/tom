// +build socket

/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
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
	proto.AssertCommandIsDefined(proto.CmdSocketShow)

	registry = append(registry, function{
		cmd:    proto.CmdSocketShow,
		handle: socketShow,
	})
}

func socketShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.SocketShow)
}

func exportSocketShow(result *proto.Result, r *msg.Result) {
	result.Socket = &[]proto.Socket{}
	*result.Socket = append(*result.Socket, r.Socket...)
}

// SocketShow function
func (m *Model) SocketShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdSocketShow,
		msg.SectionSocket,
		proto.ActionShow,
	)
	request.Socket.TomID = params.ByName(`tomID`)
	request.Socket.Namespace = r.URL.Query().Get(`namespace`)
	request.Socket.Name = r.URL.Query().Get(`name`)

	if err := request.Socket.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
		// error is ErrEmptyTomID, check for query parameter
		// supplied values
		switch {
		case request.Socket.Namespace == ``:
			fallthrough
		case request.Socket.Name == ``:
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
	m.x.Send(&w, &result, exportSocketShow)
}

// show returns full details for a specific socket
func (h *SocketReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		socketID, dictionaryID string
		tx                     *sql.Tx
		rows, links, lprops    *sql.Rows
		err                    error
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

	ct := proto.Socket{
		Namespace: q.Socket.Namespace,
		Name:      q.Socket.Name,
		Link:      []string{},
	}
	name := proto.PropertyDetail{
		Attribute: `name`,
		Value:     q.Socket.Name,
	}

	var since, until, createdAt, namedAt time.Time
	if err = txShow.QueryRow(
		q.Socket.Namespace,
		q.Socket.Name,
		txTime,
	).Scan(
		&socketID,
		&dictionaryID,
		&createdAt,
		&ct.CreatedBy,
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

	ct.CreatedAt = createdAt.Format(msg.RFC3339Milli)
	name.CreatedAt = namedAt.Format(msg.RFC3339Milli)
	name.ValidSince = since.Format(msg.RFC3339Milli)
	name.ValidUntil = until.Format(msg.RFC3339Milli)
	name.Namespace = q.Socket.Namespace
	ct.Property = make(map[string]proto.PropertyDetail)
	ct.Property[q.Socket.Namespace+`::`+ct.Name+`::name`] = name

	// fetch socket properties
	if rows, err = txProp.Query(
		dictionaryID,
		socketID,
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
		prop.Namespace = q.Socket.Namespace

		// set specialty fields
		switch prop.Attribute {
		case `name`:
			if ct.Name != prop.Value {
				rows.Close()
				mr.ExpectationFailed(
					fmt.Errorf(`Encountered confused resultset`),
				)
				return
			}
		case `type`:
			ct.Type = prop.Value
		}
		switch prop.Attribute {
		case `name`:
			// name attribute has already been added
		default:
			ct.Property[prop.Namespace+`::`+ct.Name+`::`+prop.Attribute] = prop
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch linked sockets
	linklist := [][]string{}
	if links, err = tx.Stmt(h.stmtLinked).Query(
		socketID,
		dictionaryID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for links.Next() {
		var linkedContID, linkedDictID, linkedContName, linkedDictName string
		if err = links.Scan(
			&linkedContID,
			&linkedDictID,
			&linkedContName,
			&linkedDictName,
		); err != nil {
			links.Close()
			mr.ServerError(err)
			return
		}
		ct.Link = append(ct.Link, linkedContName+`.`+linkedDictName+`.socket.tom`)
		linklist = append(linklist, []string{
			linkedContID,
			linkedDictID,
			linkedContName,
			linkedDictName,
		})
	}
	if err = links.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch properties from linked sockets
	for i := range linklist {
		if lprops, err = tx.Query(
			stmt.SocketTxShowProperties,
			linklist[i][1], // linkedDictID
			linklist[i][0], // linkedContID
			txTime,
		); err != nil {
			mr.ServerError(err)
			return
		}

		for lprops.Next() {
			prop := proto.PropertyDetail{}
			var since, until, at time.Time

			if err = lprops.Scan(
				&prop.Attribute,
				&prop.Value,
				&since,
				&until,
				&at,
				&prop.CreatedBy,
			); err != nil {
				lprops.Close()
				mr.ServerError(err)
				return
			}
			prop.ValidSince = since.Format(msg.RFC3339Milli)
			prop.ValidUntil = until.Format(msg.RFC3339Milli)
			prop.CreatedAt = at.Format(msg.RFC3339Milli)
			prop.Namespace = linklist[i][3] // linkedDictName

			// linklist[i][2] is linkedContName
			ct.Property[prop.Namespace+`::`+linklist[i][2]+`::`+prop.Attribute] = prop
		}
		if err = lprops.Err(); err != nil {
			mr.ServerError(err)
			return
		}
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Socket = append(mr.Socket, ct)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
