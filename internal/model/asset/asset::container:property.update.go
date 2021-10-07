/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset/"

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdContainerPropUpdate)

	registry = append(registry, function{
		cmd:    proto.CmdContainerPropUpdate,
		handle: containerPropUpdate,
	})
}

func containerPropUpdate(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.ContainerPropUpdate)
}

func exportContainerPropUpdate(result *proto.Result, r *msg.Result) {
	result.Container = &[]proto.Container{}
	*result.Container = append(*result.Container, r.Container...)
}

// ContainerPropUpdate function
func (m *Model) ContainerPropUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionContainer
	request.Action = proto.ActionPropUpdate

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if req.Container.Link == nil {
		req.Container.Link = []string{}
	}
	if req.Container.Property == nil {
		req.Container.Property = map[string]proto.PropertyDetail{}
	}
	request.Container = *req.Container

	request.Container.TomID = params.ByName(`tomID`)
	if err := request.Container.ParseTomID(); err != nil {
		if err != proto.ErrEmptyTomID {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
	}

	// input validation
	if err := proto.ValidNamespace(request.Container.Namespace); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if err := proto.OnlyUnreserved(request.Container.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	if request.Container.Property == nil {
		m.x.ReplyBadRequest(&w, &request, nil)
		return
	}
	for prop, obj := range request.Container.Property {
		if err := proto.OnlyUnreserved(prop); err != nil {
			m.x.ReplyBadRequest(&w, &request, err)
			return
		}
		if err := proto.CheckPropertyConstraints(&obj); err != nil {
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
	m.x.Send(&w, &result, exportContainerPropUpdate)
}

// propertyUpdate ...
func (h *ContainerWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		containerID, dictionaryID, createdAt, createdBy  string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		err                                              error
		rows                                             *sql.Rows
		tx                                               *sql.Tx
		res                                              sql.Result
		ok                                               bool
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover containerID at the start of the transaction, as the property
	// updates might include a name change
	if err = tx.QueryRow(
		stmt.ContainerTxShow,
		q.Container.Namespace,
		q.Container.Name,
		txTime,
	).Scan(
		&containerID,
		&dictionaryID,
		&createdAt,
		&createdBy,
		&nameValidSince,
		&nameValidUntil,
		&namedAt,
		&namedBy,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		tx.Rollback()
		return
	} else if err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// discover all attributes and record them with their type
	attrMap := map[string]string{}
	if rows, err = tx.Query(
		stmt.NamespaceAttributeDiscover,
		q.Container.Namespace,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	for rows.Next() {
		var attribute, typ string
		if err = rows.Scan(
			&attribute,
			&typ,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		attrMap[attribute] = typ
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// for all properties specified in the request, check that the attribute
	// exists and create missing attributes
	for key := range q.Container.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = tx.Exec(
				stmt.NamespaceAttributeAddStandard,
				q.Container.Namespace,
				key,
				q.UserIDLib,
				q.AuthUser,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
			if !mr.CheckRowsAffected(res.RowsAffected()) {
				tx.Rollback()
				return
			}
		}
	}

	// for all properties specified in the request, update the value.
	// this transparently creates missing entries.
	for key := range q.Container.Property {
		if key == `type` || q.Container.Property[key].Attribute == `type` {
			continue
		}
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Container.Property[key], containerID,
		); !ok {
			tx.Rollback()
			return
		}
		// remove the property from the map of available attributes
		delete(attrMap, q.Container.Property[key].Attribute)
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Container = append(mr.Container, q.Container)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
