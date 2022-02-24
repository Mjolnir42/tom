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

	request := msg.New(
		r, params,
		proto.CmdContainerShow,
		msg.SectionContainer,
		proto.ActionShow,
	)
	request.Container.TomID = params.ByName(`tomID`)
	request.Container.Namespace = r.URL.Query().Get(`namespace`)
	request.Container.Name = r.URL.Query().Get(`name`)
	request.Verbose, _ = strconv.ParseBool(r.URL.Query().Get(`verbose`))

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
		containerID, dictionaryID string
		tx                        *sql.Tx
		rows, links, lprops       *sql.Rows
		err                       error
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
	txResource := tx.Stmt(h.stmtTxResource)

	ct := *(proto.NewContainer())
	ct.Namespace = q.Container.Namespace
	ct.Name = q.Container.Name
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
		&containerID,
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

	if q.Verbose {
		ct.CreatedAt = createdAt.Format(msg.RFC3339Milli)
		name.CreatedAt = namedAt.Format(msg.RFC3339Milli)
		name.ValidSince = since.Format(msg.RFC3339Milli)
		name.ValidUntil = until.Format(msg.RFC3339Milli)
	} else {
		ct.CreatedBy = ``
		name.CreatedBy = ``
	}

	name.Namespace = q.Container.Namespace
	ct.Property = make(map[string]proto.PropertyDetail)
	ct.Property[q.Container.Namespace+`::`+ct.Name+`::name`] = name

	// fetch container properties
	if rows, err = txProp.Query(
		dictionaryID,
		containerID,
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
			prop.Namespace = q.Container.Namespace
		} else {
			prop.CreatedBy = ``
		}

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

	// query parent information
	var rteID, rteDictID, rteDictName, rteName string
	noParent := false
	if err = tx.Stmt(
		h.stmtTxParent,
	).QueryRow(
		containerID,
		txTime,
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
		ct.Parent = (&proto.Runtime{
			ID:        rteID,
			Namespace: rteDictName,
			Name:      rteName,
		}).FormatTomID()
	}

	// fetch resource links for current containerID
	var resource string
	noResource := false
	if err = txResource.QueryRow(
		q.Container.Namespace,
		containerID,
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
		ct.Resources = append(ct.Resources, resource)
	}

	// fetch linked containers
	linklist := [][]string{}
	if links, err = tx.Stmt(h.stmtLinked).Query(
		containerID,
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
		ct.Link = append(ct.Link, linkedContName+`.`+linkedDictName+`.container.tom`)
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
			ct.Resources = append(ct.Resources, linkResource)
		}
	}

	// fetch properties from linked containers
	for i := range linklist {
		if lprops, err = tx.Query(
			stmt.ContainerTxShowProperties,
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
			if q.Verbose {
				prop.ValidSince = since.Format(msg.RFC3339Milli)
				prop.ValidUntil = until.Format(msg.RFC3339Milli)
				prop.CreatedAt = at.Format(msg.RFC3339Milli)
			} else {
				prop.CreatedBy = ``
			}
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
	mr.Container = append(mr.Container, ct)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
