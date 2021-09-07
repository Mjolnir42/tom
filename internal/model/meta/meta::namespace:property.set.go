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
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespacePropSet)

	registry = append(registry, function{
		cmd:    proto.CmdNamespacePropSet,
		handle: namespacePropertySet,
	})
}

func namespacePropertySet(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespacePropertySet)
}

func exportNamespacePropertySet(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespacePropertySet function
func (m *Model) NamespacePropertySet(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = msg.ActionPropSet

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

	if err := proto.OnlyUnreserved(request.Namespace.Name); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	// check property structure is setup
	if request.Namespace.Property == nil {
		m.x.ReplyBadRequest(&w, &request, fmt.Errorf(
			"Invalid property.set request without properties",
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

// propertySet sets the properties of a namespace to the exact set of
// properties configured in the request.
// It creates, updates and deactivates properties as required.
func (h *NamespaceWriteHandler) propertySet(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		tx       *sql.Tx
		txTime   time.Time
		ok, done bool
		rows     *sql.Rows
		res      sql.Result
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

	// for all properties specified in the request, check that the attribute
	// exists and create missing attributes
	for key := range q.Namespace.Property {
		if _, ok = attrMap[key]; !ok {
			if res, err = tx.Stmt(h.stmtAttStdAdd).Exec(
				q.Namespace.Name,
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
		// add newly created attribute to map
		attrMap[key] = proto.AttributeStandard
	}

	// for all properties specified in the request, update the value.
	// this transparently creates missing entries.
	for key := range q.Namespace.Property {
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Namespace.Property[key],
		); !ok {
			tx.Rollback()
			return
		}
		// remove the property from the map of available attributes
		delete(attrMap, q.Namespace.Property[key].Attribute)
	}

	// special handling here for namespaces: do not clear unspecified
	// attributes with `dict_` prefix. Remove them from the map of attributes
	// that require clearing
	for key := range attrMap {
		if strings.HasPrefix(key, `dict_`) {
			delete(attrMap, key)
		}
	}

	// for all properties not in the request, any currently valid value must
	// be invalided by setting its validUntil value to the time of the
	// transaction
	for key := range attrMap {
		var stmtSelect, stmtClamp *sql.Stmt

		switch attrMap[key] {
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
