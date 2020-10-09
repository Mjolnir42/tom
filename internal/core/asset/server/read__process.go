/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package server // import "github.com/mjolnir42/tom/internal/core/asset/server"

import (
	"database/sql"
	"fmt"

	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

// process is the request dispatcher
func (h *ReadHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)
	//	logRequest(h.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		h.list(q, &result)
	case msg.ActionShow:
		h.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all servers
func (h *ReadHandler) list(q *msg.Request, mr *msg.Result) {
	var (
		id, namespace, name, typ string
		rows                     *sql.Rows
		err                      error
	)

	if rows, err = h.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&namespace,
			&name,
			&typ,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		mr.Server = append(mr.Server, proto.Server{
			ID:        id,
			Namespace: namespace,
			Name:      name,
			Type:      typ,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// show returns full details for a specific server
func (h *ReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		tx                                           *sql.Tx
		err                                          error
		txFind, txAttr, txParent, txLink             *sql.Stmt
		qrySrvID, qrySrvName, qryDictID, qryDictName *sql.NullString
		rows                                         *sql.Rows
		server                                       proto.Server
		ambiguous                                    bool
		id, dictID, dictName, attrID, key, value     string
		rteID, rteDictID, rteDictName, rteName       string
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

	// find the server
	if qrySrvID.String = q.Server.ID; qrySrvID.String != `` {
		qrySrvID.Valid = true
	}
	if qrySrvName.String = q.Server.Name; qrySrvName.String != `` {
		qrySrvName.Valid = true
	}
	if qryDictName.String = q.Server.Namespace; qryDictName.String != `` {
		qryDictName.Valid = true
	}
	txFind = tx.Stmt(h.stmtFind)
	if rows, err = txFind.Query(
		qrySrvName,
		qrySrvID,
		qryDictID,
		qryDictName,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for rows.Next() {
		if ambiguous {
			rows.Close()
			mr.ExpectationFailed(fmt.Errorf(`Request is ambiguous`))
			return
		}
		if err = rows.Scan(
			&id,
			&dictID,
			&dictName,
			&attrID,
			&key,
			&value,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		server = proto.Server{
			ID:        id,
			Namespace: dictName,
			Name:      value,
		}
		ambiguous = true
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// query all server attributes
	txAttr = tx.Stmt(h.stmtAttribute)
	if rows, err = txAttr.Query(
		server.ID,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(
			&id,
			&dictName,
			&key,
			&value,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		switch {
		case id != server.ID || dictName != server.Namespace:
			rows.Close()
			mr.ExpectationFailed(fmt.Errorf(`Request is ambiguous`))
			return
		case key == `type`:
			server.Type = value
		case key == `name`:
			server.Name = value
		default:
			server.Property = append(server.Property, [2]string{key, value})
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	// query parent
	txParent = tx.Stmt(h.stmtParent)
	if err = txParent.QueryRow(
		server.ID,
	).Scan(
		&rteID,
		&rteDictID,
		&rteDictName,
		&rteName,
	); err == sql.ErrNoRows {
		// not an error
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	server.Parent = (&proto.Runtime{
		ID:        rteID,
		Namespace: rteDictName,
		Name:      rteName,
	}).TomID()

	// query links
	qrySrvName.String = ``
	qrySrvName.Valid = false
	qryDictName.String = ``
	qryDictName.Valid = false
	txLink = tx.Stmt(h.stmtLink)
	if rows, err = txLink.Query(
		server.ID,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(
			&id,
			&dictID,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}

		if qrySrvID.String = id; qrySrvID.String != `` {
			qrySrvID.Valid = true
		}
		if qryDictID.String = dictID; qryDictID.String != `` {
			qryDictID.Valid = true
		}
		if err = tx.Stmt(h.stmtFind).QueryRow(
			qrySrvName,
			qrySrvID,
			qryDictID,
			qryDictName,
		).Scan(
			id,
			dictID,
			dictName,
			attrID,
			key,
			value,
		); err == sql.ErrNoRows {
			rows.Close()
			mr.ServerError(fmt.Errorf(`Inconsistent dataset`))
			return
		} else if err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		server.Link = append(server.Link, (&proto.Server{
			Namespace: dictName,
			Name:      value,
		}).TomID())
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
