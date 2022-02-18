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
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespacePropRemove)

	registry = append(registry, function{
		cmd:    proto.CmdNamespacePropRemove,
		handle: namespacePropertyRemove,
	})
}

func namespacePropertyRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespacePropertyRemove)
}

func exportNamespacePropertyRemove(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespacePropertyRemove function
func (m *Model) NamespacePropertyRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(
		r, params,
		proto.CmdNamespacePropRemove,
		msg.SectionNamespace,
		proto.ActionPropRemove,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = *req.Namespace
	request.Namespace.TomID = params.ByName(`tomID`)
	if err := request.Namespace.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if err := proto.ValidNamespace(request.Namespace.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	// check property structure is setup
	if request.Namespace.Property == nil {
		m.x.ReplyBadRequest(&w, &request, fmt.Errorf(
			"Invalid %s request without properties",
			proto.ActionPropRemove,
		))
		return
	}

	// check property names are all valid
	for prop := range request.Namespace.Property {
		if err := proto.OnlyUnreserved(prop); err != nil {
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
	m.x.Send(&w, &result, exportNamespacePropertySet)
}

// propertyRemove deactivates properties of a namespace:
func (h *NamespaceWriteHandler) propertyRemove(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		tx       *sql.Tx
		txTime   time.Time
		ok, done bool
		rows     *sql.Rows
	)
	// setup a consistent transaction time timestamp that is used for all
	// record updates during the transaction to make them concurrent
	txTime = time.Now().UTC()

	// tx.Begin()
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover all attributes and record them with their type
	attrMap := map[string]string{}
	reqMap := map[string]string{}
	if rows, err = tx.Stmt(h.stmtAttDiscover).Query(
		q.Namespace.Name,
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

	// special handling, do not allow to remove the name
	// or type of the namespace
	delete(q.Namespace.Property, `dict_name`)
	delete(q.Namespace.Property, `dict_type`)

	// for all properties specified in the request, check that the attribute
	// exists and record its type
	for key := range q.Namespace.Property {
		if _, ok = attrMap[q.Namespace.Property[key].Attribute]; ok {
			reqMap[q.Namespace.Property[key].Attribute] = attrMap[q.Namespace.Property[key].Attribute]
		} else {
			mr.BadRequest(fmt.Errorf(
				"Specified property attribute <%s> does not exist.",
				key,
			))
			tx.Rollback()
			return
		}
	}

	// for all properties in the request, any currently valid value must
	// be invalided by setting its validUntil value to the time of the
	// transaction
	for key := range reqMap {
		var stmtSelect, stmtClamp *sql.Stmt

		switch reqMap[key] {
		case proto.AttributeStandard:
			stmtSelect = tx.Stmt(h.stmtTxStdPropSelect)
			stmtClamp = tx.Stmt(h.stmtTxStdPropClamp)
		case proto.AttributeUnique:
			stmtSelect = tx.Stmt(h.stmtTxUniqPropSelect)
			stmtClamp = tx.Stmt(h.stmtTxUniqPropClamp)
		default:
			mr.ServerError()
			tx.Rollback()
			return
		}

		if ok, done = h.txPropClamp(
			q, mr, tx, &txTime, stmtSelect, stmtClamp,
			proto.PropertyDetail{
				Attribute: key,
				// construct an imaginary new value for the property. The clamping
				// function does not invalidate the current value, if the provided
				// imaginary record matches the existing record in both value and
				// upper validUntil bound
				Value: txTime.Format(msg.RFC3339Milli) + key + `_clamp`,
			},
			// send the transaction time as the imaginary new value's requested
			// validSince:
			// - guaranteed to be after the current value's validSince, as it
			//   otherwise would not be valid
			// - this value is used as the new validUntil boundary of the exitsing
			//   value
			txTime,
			// send the transaction time as the imaginary new value's requested
			// validUntil. This is matched against the existing value's validUntil
			// to determine if the record needs clamping
			txTime,
		); !ok {
			// appropriate error function is already called by h.txPropClamp
			tx.Rollback()
			return
		} else if done {
			// txPropClamp did nothing because the existing record matched both
			// time and value. Retry with a different value.
			if ok, done = h.txPropClamp(
				q, mr, tx, &txTime, stmtSelect, stmtClamp,
				proto.PropertyDetail{
					Attribute: key,
					Value:     txTime.Format(msg.RFC3339Milli) + key + `_alternate`,
				},
				txTime,
				txTime,
			); !ok {
				tx.Rollback()
				return
			} else if done {
				// this should not be possible -> abort.
				mr.ServerError(fmt.Errorf(
					"txPropClamp encountered impossible repeat matches for %s",
					key,
				))
				tx.Rollback()
				return
			}
		}
	}

	// tx.Commit()
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
