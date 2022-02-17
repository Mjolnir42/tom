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
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdRuntimeShow)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeShow,
		handle: runtimeShow,
	})
}

func runtimeShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeShow)
}

func exportRuntimeShow(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimeShow function
func (m *Model) RuntimeShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdRuntimeShow,
		msg.SectionRuntime,
		proto.ActionShow,
	)
	request.Runtime.TomID = params.ByName(`tomID`)
	request.Runtime.Namespace = r.URL.Query().Get(`namespace`)
	request.Runtime.Name = r.URL.Query().Get(`name`)

	if err := request.Runtime.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
		// error is ErrEmptyTomID, check for query parameter
		// supplied values
		switch {
		case request.Runtime.Namespace == ``:
			fallthrough
		case request.Runtime.Name == ``:
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
	m.x.Send(&w, &result, exportRuntimeShow)
}

// show returns full details for a specific runtime
func (h *RuntimeReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		rteID, dictionaryID string
		tx                  *sql.Tx
		rows, links, lprops *sql.Rows
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
	txParent := tx.Stmt(h.stmtParent)
	txChildren := tx.Stmt(h.stmtTxChildren)
	txResource := tx.Stmt(h.stmtTxResource)

	rte := proto.Runtime{
		Namespace: q.Runtime.Namespace,
		Name:      q.Runtime.Name,
		Link:      []string{},
		Children:  []string{},
	}
	name := proto.PropertyDetail{
		Attribute: `name`,
		Value:     q.Runtime.Name,
	}

	var since, until, createdAt, namedAt time.Time
	if err = txShow.QueryRow(
		q.Runtime.Namespace,
		q.Runtime.Name,
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
	name.Namespace = q.Runtime.Namespace
	rte.Property = make(map[string]proto.PropertyDetail)
	rte.Property[q.Runtime.Namespace+`::`+rte.Name+`::name`] = name

	// fetch runtime properties
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
		prop.ValidSince = since.Format(msg.RFC3339Milli)
		prop.ValidUntil = until.Format(msg.RFC3339Milli)
		prop.CreatedAt = at.Format(msg.RFC3339Milli)
		prop.Namespace = q.Runtime.Namespace

		switch {
		case strings.HasSuffix(prop.Attribute, `_json`):
			fallthrough
		case strings.HasSuffix(prop.Attribute, `_list`):
			prop.Raw = []byte(prop.Value)
		}

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
			rte.Property[prop.Namespace+`::`+rte.Name+`::`+prop.Attribute] = prop
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch parent information for stacked runtimes
	var hasParent bool = true
	var ptEntity, ptObjID, ptDictID, ptDictName, ptObjName string
	if err = txParent.QueryRow(
		rteID,
		txTime,
	).Scan(
		&ptEntity,
		&ptObjID,
		&ptDictID,
		&ptDictName,
		&ptObjName,
	); err == sql.ErrNoRows {
		// not an error
		hasParent = false
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	if hasParent {
		switch ptEntity {
		case proto.EntityRuntime:
			rte.Parent = (&proto.Runtime{
				Namespace: ptDictName,
				Name:      ptObjName,
			}).FormatTomID()
		case proto.EntityServer:
			rte.Parent = (&proto.Server{
				Namespace: ptDictName,
				Name:      ptObjName,
			}).FormatTomID()
		case proto.EntityOrchestration:
			rte.Parent = (&proto.Orchestration{
				Namespace: ptDictName,
				Name:      ptObjName,
			}).FormatTomID()
		default:
			mr.ServerError()
			return
		}
	}

	// query children
	if rows, err = txChildren.Query(
		rteID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		var chldEntity, chldName, chldDictName string
		if err = rows.Scan(
			&chldEntity,
			&chldName,
			&chldDictName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			return
		}
		switch chldEntity {
		case proto.EntityRuntime:
			rte.Children = append(rte.Children, (&proto.Runtime{
				Namespace: chldDictName,
				Name:      chldName,
			}).FormatTomID())
		case proto.EntityServer:
			rte.Children = append(rte.Children, (&proto.Server{
				Namespace: chldDictName,
				Name:      chldName,
			}).FormatTomID())
		case proto.EntityOrchestration:
			rte.Children = append(rte.Children, (&proto.Orchestration{
				Namespace: chldDictName,
				Name:      chldName,
			}).FormatTomID())
		default:
			rows.Close()
			mr.ServerError()
			return
		}
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch resource links for current rteID
	var resource string
	noResource := false
	if err = txResource.QueryRow(
		q.Runtime.Namespace,
		rteID,
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
		rte.Resources = append(rte.Resources, resource)
	}

	// fetch linked runtimes
	linklist := [][]string{}
	if links, err = tx.Stmt(h.stmtLinked).Query(
		rteID,
		dictionaryID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for links.Next() {
		var linkedRteID, linkedDictID, linkedRteName, linkedDictName string
		if err = links.Scan(
			&linkedRteID,
			&linkedDictID,
			&linkedRteName,
			&linkedDictName,
		); err != nil {
			links.Close()
			mr.ServerError(err)
			return
		}

		rte.Link = append(rte.Link, (&proto.Runtime{
			Namespace: linkedDictName,
			Name:      linkedRteName,
		}).FormatTomID())

		linklist = append(linklist, []string{
			linkedRteID,
			linkedDictID,
			linkedRteName,
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
			rte.Resources = append(rte.Resources, linkResource)
		}
	}

	for i := range linklist {
		// fetch properties from linked runtime
		if lprops, err = tx.Query(
			stmt.RuntimeTxShowProperties,
			linklist[i][1], // linkedDictID
			linklist[i][0], // linkedRteID
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

			switch {
			case strings.HasSuffix(prop.Attribute, `_json`):
				fallthrough
			case strings.HasSuffix(prop.Attribute, `_list`):
				prop.Raw = []byte(prop.Value)
			}

			// linklist[i][2] is linkedRteName
			rte.Property[prop.Namespace+`::`+linklist[i][2]+`::`+prop.Attribute] = prop
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
	mr.Runtime = append(mr.Runtime, rte)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
