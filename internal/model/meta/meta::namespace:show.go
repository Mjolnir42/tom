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
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespaceShow)

	registry = append(registry, function{
		cmd:    proto.CmdNamespaceShow,
		handle: namespaceShow,
	})
}

func namespaceShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespaceShow)
}

func exportNamespaceShow(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespaceShow function
func (m *Model) NamespaceShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = proto.ActionShow
	request.Namespace = proto.Namespace{
		TomID: params.ByName(`tomID`),
		Name:  r.URL.Query().Get(`name`),
	}

	if err := request.Namespace.ParseTomID(); err != nil {
		if !(err == proto.ErrEmptyTomID && request.Namespace.Name != ``) {
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
	m.x.Send(&w, &result, exportNamespaceShow)
}

// show returns full details for a specific server
func (h *NamespaceReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		dictionaryID string
		tx           *sql.Tx
		rows         *sql.Rows
		err          error
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
	txAttr := tx.Stmt(h.stmtAttr)

	ns := proto.Namespace{}

	if err = txShow.QueryRow(
		// unique constraint on dictionary name
		q.Namespace.Name,
	).Scan(
		&dictionaryID,
		&ns.Name,
		&ns.CreatedAt,
		&ns.CreatedBy,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	ns.Property = make(map[string]proto.PropertyDetail)
	if rows, err = txProp.Query(
		dictionaryID,
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
		prop.Namespace = q.Namespace.Name
		ns.Property[prop.Attribute] = prop

		// set sepcialty fields for well known namespace properties
		switch prop.Attribute {
		case `dict_name`:
			if ns.Name != prop.Value {
				rows.Close()
				mr.ExpectationFailed(
					fmt.Errorf(`Encountered confused resultset`),
				)
				return
			}
		case `dict_type`:
			ns.Type = prop.Value
		case `dict_lookup`:
			ns.LookupKey = prop.Value
		case `dict_uri`:
			ns.LookupURI = prop.Value
		case `dict_ntt_list`:
			ns.Constraint = strings.Split(prop.Value, `,`)
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	ns.Attributes = []proto.AttributeDefinition{}
	if rows, err = txAttr.Query(
		dictionaryID,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		attr := proto.AttributeDefinition{}
		var at time.Time
		var author, typ string

		if err = rows.Scan(
			&attr.Key,
			&typ,
			&at,
			&author,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}

		switch typ {
		case proto.AttributeUnique:
			attr.Unique = true
		}
		ns.Attributes = append(ns.Attributes, attr)
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
	mr.Namespace = append(mr.Namespace, ns)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
