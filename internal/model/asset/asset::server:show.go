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
	"strconv"
	"strings"
	"time"

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

	request := msg.New(
		r, params,
		proto.CmdServerShow,
		msg.SectionServer,
		proto.ActionShow,
	)
	request.Server.TomID = params.ByName(`tomID`)
	request.Server.Namespace = r.URL.Query().Get(`namespace`)
	request.Server.Name = r.URL.Query().Get(`name`)
	request.Verbose, _ = strconv.ParseBool(r.URL.Query().Get(`verbose`))

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
		serverID, dictionaryID                 string
		tx                                     *sql.Tx
		rows, links, lprops                    *sql.Rows
		err                                    error
		rteID, rteDictID, rteDictName, rteName string
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
	txShow := tx.Stmt(h.stmtTxShow)
	txProp := tx.Stmt(h.stmtTxProp)
	txChildren := tx.Stmt(h.stmtTxChildren)
	txResource := tx.Stmt(h.stmtTxResource)

	server := proto.NewServer()
	server.Namespace = q.Server.Namespace
	server.Name = q.Server.Name
	name := proto.PropertyDetail{
		Attribute: `name`,
		Value:     q.Server.Name,
	}

	var since, until, createdAt, namedAt time.Time
	if err = txShow.QueryRow(
		q.Server.Namespace,
		q.Server.Name,
		txTime,
	).Scan(
		&serverID,
		&dictionaryID,
		&createdAt,
		&server.CreatedBy,
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

	if q.Verbose {
		server.CreatedAt = createdAt.Format(msg.RFC3339Milli)
		name.CreatedAt = namedAt.Format(msg.RFC3339Milli)
		name.ValidSince = since.Format(msg.RFC3339Milli)
		name.ValidUntil = until.Format(msg.RFC3339Milli)
	} else {
		server.CreatedBy = ``
		name.CreatedBy = ``
	}
	name.Namespace = q.Server.Namespace
	server.Property = make(map[string]proto.PropertyDetail)
	server.Property[q.Server.Namespace+`::`+server.Name+`::name`] = name

	// fetch server properties
	if rows, err = txProp.Query(
		dictionaryID,
		serverID,
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
		if q.Verbose {
			prop.ValidSince = since.Format(msg.RFC3339Milli)
			prop.ValidUntil = until.Format(msg.RFC3339Milli)
			prop.CreatedAt = at.Format(msg.RFC3339Milli)
		} else {
			prop.CreatedBy = ``
		}
		prop.Namespace = q.Server.Namespace

		switch {
		case strings.HasSuffix(prop.Attribute, `_json`):
			fallthrough
		case strings.HasSuffix(prop.Attribute, `_list`):
			prop.Raw = []byte(prop.Value)
		}

		// set specialty fields
		switch prop.Attribute {
		case `name`:
			if server.Name != prop.Value {
				rows.Close()
				mr.ExpectationFailed(
					fmt.Errorf(`Encountered confused resultset`),
				)
				return
			}
		case `type`:
			server.Type = prop.Value
		}
		switch prop.Attribute {
		case `name`:
			// name attribute has already been added
		default:
			server.Property[prop.Namespace+`::`+server.Name+`::`+prop.Attribute] = prop
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// query parent
	noParent := false
	if err = tx.Stmt(
		h.stmtParent,
	).QueryRow(
		serverID,
	).Scan(
		&rteID,
		&rteDictID,
		&rteDictName,
		&rteName,
	); err == sql.ErrNoRows {
		// not an error
		noParent = true
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	if !noParent {
		server.Parent = (&proto.Runtime{
			ID:        rteID,
			Namespace: rteDictName,
			Name:      rteName,
		}).FormatTomID()
	}

	// query children
	if rows, err = txChildren.Query(
		serverID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		var chldName, chldDictName string
		if err = rows.Scan(
			&chldName,
			&chldDictName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		// servers can only have runtimes as children
		server.Children = append(server.Children, (&proto.Runtime{
			Namespace: chldDictName,
			Name:      chldName,
		}).FormatTomID())
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch resource links for current serverID
	var resource string
	noResource := false
	if err = txResource.QueryRow(
		q.Server.Namespace,
		serverID,
		txTime,
	).Scan(
		&resource,
	); err == sql.ErrNoRows {
		// not an error, might not be a referential namespace
		noResource = true
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	if !noResource {
		server.Resources = append(server.Resources, resource)
	}

	// fetch linked servers
	linklist := [][]string{}
	if links, err = tx.Stmt(h.stmtLinked).Query(
		serverID,
		dictionaryID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for links.Next() {
		var linkedSrvID, linkedDictID, linkedSrvName, linkedDictName string
		if err = links.Scan(
			&linkedSrvID,
			&linkedDictID,
			&linkedSrvName,
			&linkedDictName,
		); err != nil {
			links.Close()
			mr.ServerError(err)
			return
		}

		server.Link = append(server.Link, (&proto.Server{
			ID:        linkedSrvID,
			Namespace: linkedDictName,
			Name:      linkedSrvName,
		}).FormatTomID())

		linklist = append(linklist, []string{
			linkedSrvID,
			linkedDictID,
			linkedSrvName,
			linkedDictName,
		})
	}
	if err = links.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch linked resources
	for i := range linklist {
		noResource = false
		var linkResource string
		if err = txResource.QueryRow(
			linklist[i][3],
			linklist[i][0],
			txTime,
		).Scan(
			&linkResource,
		); err == sql.ErrNoRows {
			// not an error, might not be a referential namespace
			noResource = true
		} else if err != nil {
			mr.ServerError(err)
			return
		}
		if !noResource {
			server.Resources = append(server.Resources, linkResource)
		}
	}

	for i := range linklist {
		// fetch properties from linked runtime
		if lprops, err = tx.Stmt(
			h.stmtTxProp,
		).Query(
			linklist[i][1], // linkedDictID
			linklist[i][0], // linkedSrvID
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
			if q.Verbose {
				prop.ValidSince = since.Format(msg.RFC3339Milli)
				prop.ValidUntil = until.Format(msg.RFC3339Milli)
				prop.CreatedAt = at.Format(msg.RFC3339Milli)
			} else {
				prop.CreatedBy = ``
			}
			prop.Namespace = linklist[i][3] // linkedDictName

			switch {
			case strings.HasSuffix(prop.Attribute, `_json`):
				fallthrough
			case strings.HasSuffix(prop.Attribute, `_list`):
				prop.Raw = []byte(prop.Value)
			}

			// linklist[i][2] is linkedSrvName
			server.Property[prop.Namespace+`::`+linklist[i][2]+`::`+prop.Attribute] = prop
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
	mr.Server = append(mr.Server, *server)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
