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

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdServerShow)

	registry = append(registry, function{
		cmd:    proto.CmdServerShow,
		handle: serverShow,
	})
}

func serverShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ServerShow)
}

func exportServerShow(result *proto.Result, r *msg.Result) {
	result.Server = &[]proto.Server{}
	*result.Server = append(*result.Server, r.Server...)
}

// ServerShow function
func (m *Model) ServerShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionServer
	request.Action = proto.ActionShow
	request.Server = proto.Server{
		TomID:     params.ByName(`tomID`),
		Namespace: r.URL.Query().Get(`namespace`),
		Name:      r.URL.Query().Get(`name`),
	}

	if err := request.Server.ParseTomID(); err != nil {
		if !(err == proto.ErrEmptyTomID && request.Server.Name != ``) {
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
	m.x.Send(&w, &result, exportServerShow)
}

// show returns full details for a specific server
func (h *ServerReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		tx                                           *sql.Tx
		err                                          error
		txFind, txAttr, txParent, txLink             *sql.Stmt
		qrySrvID, qrySrvName, qryDictID, qryDictName *sql.NullString
		rows                                         *sql.Rows
		server                                       proto.Server
		ambiguous                                    bool
		id, dictID, dictName, attrID, attrTyp, value string
		rteID, rteDictID, rteDictName, rteName, key  string
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
			&attrTyp,
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
			server.Property[key] = proto.PropertyDetail{
				Attribute: key,
				Value:     value,
			}
			switch attrTyp {
			case proto.AttributeUnique:
				server.UniqProperty = append(server.UniqProperty, server.Property[key])
			case proto.AttributeStandard:
				server.StdProperty = append(server.StdProperty, server.Property[key])
			default:
				rows.Close()
				mr.ExpectationFailed(fmt.Errorf(`Received impossible attribute type`))
				return
			}
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
	}).FormatTomID()

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
		}).FormatTomID())
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
