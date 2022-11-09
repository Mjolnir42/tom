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
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdOrchestrationShow)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationShow,
		handle: orchestrationShow,
	})
}

func orchestrationShow(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationShow)
}

func exportOrchestrationShow(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationShow function
func (m *Model) OrchestrationShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdOrchestrationShow,
		msg.SectionOrchestration,
		proto.ActionShow,
	)
	request.Orchestration.TomID = params.ByName(`tomID`)
	request.Orchestration.Namespace = r.URL.Query().Get(`namespace`)
	request.Orchestration.Name = r.URL.Query().Get(`name`)

	if err := request.Orchestration.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
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
	m.x.Send(&w, &result, exportOrchestrationShow)
}

// show returns full details about an orchestration
func (h *OrchestrationReadHandler) show(q *msg.Request, mr *msg.Result) {
	var (
		oreID, namespaceID              string
		tx                              *sql.Tx
		props, links, parents, children *sql.Rows
		err                             error
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
	txLinks := tx.Stmt(h.stmtTxLinks)
	txProp := tx.Stmt(h.stmtTxProp)
	txParent := tx.Stmt(h.stmtTxParent)
	txChildren := tx.Stmt(h.stmtTxChildren)
	txResource := tx.Stmt(h.stmtTxResource)

	ore := *(proto.NewOrchestration())
	ore.Namespace = q.Orchestration.Namespace
	ore.Name = q.Orchestration.Name
	name := proto.PropertyDetail{
		Attribute: `name`,
		Value:     q.Orchestration.Name,
	}

	// fetch base orchestration
	var since, until, createdAt, namedAt time.Time
	if err = txShow.QueryRow(
		q.Orchestration.Namespace,
		q.Orchestration.Name,
		txTime,
	).Scan(
		&oreID,
		&namespaceID,
		&createdAt,
		&ore.CreatedBy,
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

	ore.CreatedAt = createdAt.Format(msg.RFC3339Milli)
	name.CreatedAt = namedAt.Format(msg.RFC3339Milli)
	name.ValidSince = since.Format(msg.RFC3339Milli)
	name.ValidUntil = until.Format(msg.RFC3339Milli)
	name.Namespace = q.Orchestration.Namespace
	ore.Property = make(map[string]proto.PropertyDetail)
	ore.Property[q.Orchestration.Namespace+`::`+ore.Name+`::name`] = name

	// fetch orchestration properties
	if props, err = txProp.Query(
		namespaceID,
		oreID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for props.Next() {
		prop := proto.PropertyDetail{
			Namespace: q.Orchestration.Namespace,
		}
		var since, until, at time.Time

		if err = props.Scan(
			&prop.Attribute,
			&prop.Value,
			&since,
			&until,
			&at,
			&prop.CreatedBy,
		); err != nil {
			props.Close()
			mr.ServerError(err)
			return
		}
		prop.ValidSince = since.Format(msg.RFC3339Milli)
		prop.ValidUntil = until.Format(msg.RFC3339Milli)
		prop.CreatedAt = at.Format(msg.RFC3339Milli)

		// set structured fields
		switch {
		case strings.HasSuffix(prop.Attribute, `_json`):
			fallthrough
		case strings.HasSuffix(prop.Attribute, `_list`):
			prop.Raw = []byte(prop.Value)
		}

		// set specialty fields
		switch prop.Attribute {
		case `name`:
			if ore.Name != prop.Value {
				props.Close()
				mr.ExpectationFailed(
					fmt.Errorf(`Encountered confused resultset`),
				)
				return
			}
		case `type`:
			ore.Type = prop.Value
		}

		// copy property for export
		switch prop.Attribute {
		case `name`:
			// name attribute has already been added
		default:
			ore.Property[prop.Namespace+`::`+ore.Name+`::`+prop.Attribute] = prop
		}
	}
	if err = props.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch resource links for current oreID
	var resource string
	noResource := false
	if err = txResource.QueryRow(
		q.Orchestration.Namespace,
		oreID,
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
		ore.Resources = append(ore.Resources, resource)
	}

	// fetch linked orchestration environments
	linklist := [][]string{}
	if links, err = txLinks.Query(
		oreID,
		namespaceID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for links.Next() {
		var linkedOreID, linkedNsID, linkedOreName, linkedNsName string
		if err = links.Scan(
			&linkedOreID,
			&linkedNsID,
			&linkedOreName,
			&linkedNsName,
		); err != nil {
			links.Close()
			mr.ServerError(err)
			return
		}

		ore.Link = append(ore.Link, (&proto.Orchestration{
			Namespace: linkedNsName,
			Name:      linkedOreName,
		}).FormatTomID())

		linklist = append(linklist, []string{
			linkedOreID,
			linkedNsID,
			linkedOreName,
			linkedNsName,
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
			ore.Resources = append(ore.Resources, linkResource)
		}
	}

	// fetch all properties from linked environments
	for i := range linklist {
		if props, err = txProp.Query(
			linklist[i][1], // linkedDictID
			linklist[i][0], // linkedOreID
			txTime,
		); err != nil {
			mr.ServerError(err)
			return
		}

		for props.Next() {
			prop := proto.PropertyDetail{}
			var since, until, at time.Time

			if err = props.Scan(
				&prop.Attribute,
				&prop.Value,
				&since,
				&until,
				&at,
				&prop.CreatedBy,
			); err != nil {
				props.Close()
				mr.ServerError(err)
				return
			}
			prop.ValidSince = since.Format(msg.RFC3339Milli)
			prop.ValidUntil = until.Format(msg.RFC3339Milli)
			prop.CreatedAt = at.Format(msg.RFC3339Milli)
			prop.Namespace = linklist[i][3] // linkedNsName

			switch {
			case strings.HasSuffix(prop.Attribute, `_json`):
				fallthrough
			case strings.HasSuffix(prop.Attribute, `_list`):
				prop.Raw = []byte(prop.Value)
			}

			// linklist[i][2] is linkedOreName
			ore.Property[prop.Namespace+`::`+linklist[i][2]+`::`+prop.Attribute] = prop
		}
		if err = props.Err(); err != nil {
			mr.ServerError(err)
			return
		}
	}

	// fetch parent information for stacked orchestration environments
	var ptObjID, ptDictID, ptDictName, ptObjName string
	if parents, err = txParent.Query(
		oreID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}
	for parents.Next() {
		if err = parents.Scan(
			&ptObjID,
			&ptDictID,
			&ptDictName,
			&ptObjName,
		); err != nil {
			parents.Close()
			mr.ServerError(err)
			return
		}

		ore.Parent = append(ore.Parent, (&proto.Runtime{
			Namespace: ptDictName,
			Name:      ptObjName,
		}).FormatTomID())
	}
	if err = parents.Err(); err != nil {
		mr.ServerError(err)
		return
	}

	// fetch children
	if children, err = txChildren.Query(
		oreID,
		txTime,
	); err != nil {
		mr.ServerError(err)
		return
	}

	for children.Next() {
		var childName, childNamespace string
		if err = children.Scan(
			&childName,
			&childNamespace,
		); err != nil {
			children.Close()
			mr.ServerError(err)
			return
		}

		ore.Children = append(ore.Children, (&proto.Runtime{
			Namespace: childNamespace,
			Name:      childName,
		}).FormatTomID())
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Orchestration = append(mr.Orchestration, ore)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
