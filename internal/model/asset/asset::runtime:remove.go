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
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdRuntimeRemove)

	registry = append(registry, function{
		cmd:    proto.CmdRuntimeRemove,
		handle: runtimeRemove,
	})
}

func runtimeRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.RuntimeRemove)
}

func exportRuntimeRemove(result *proto.Result, r *msg.Result) {
	result.Runtime = &[]proto.Runtime{}
	*result.Runtime = append(*result.Runtime, r.Runtime...)
}

// RuntimeRemove function
func (m *Model) RuntimeRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdRuntimeRemove,
		msg.SectionRuntime,
		proto.ActionRemove,
	)
	request.Runtime.TomID = params.ByName(`tomID`)

	if err := request.Runtime.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportRuntimeRemove)
}

// remove invalidates an existing runtime
func (h *RuntimeWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		rteID, dictionaryID, createdAt, createdBy        string
		nameValidSince, nameValidUntil, namedAt, namedBy string
		err                                              error
		rows                                             *sql.Rows
		tx                                               *sql.Tx
		ok, done                                         bool
	)
	// setup a consistent transaction time timestamp that is used for all
	// records
	txTime = time.Now().UTC()

	// open transaction
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover rteID at the start of the transaction
	if err = tx.QueryRow(
		stmt.RuntimeTxShow,
		q.Runtime.Namespace,
		q.Runtime.Name,
		txTime,
	).Scan(
		&rteID,
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

	// unstack the runtime itself
	txStackClamp := tx.Stmt(h.stmtTxStackClamp)
	if _, err = txStackClamp.Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// unstack all runtime children currently on the runtime
	if _, err = tx.Stmt(
		h.stmtTxCldRte,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Stmt(
		h.stmtTxCldRteClean,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// unstack all server children currently on the runtime
	if _, err = tx.Stmt(
		h.stmtTxCldSrv,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Stmt(
		h.stmtTxCldSrvClean,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// unstack all socket children currently on the runtime
	if _, err = tx.Stmt(
		h.stmtTxCldSok,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Stmt(
		h.stmtTxCldSokClean,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// unstack all container children currently on the runtime
	if _, err = tx.Stmt(
		h.stmtTxCldCnr,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Stmt(
		h.stmtTxCldCnrClean,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// unstack all orchestration children currently on the runtime
	if _, err = tx.Stmt(
		h.stmtTxCldOre,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Stmt(
		h.stmtTxCldOreClean,
	).Exec(
		txTime,
		rteID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// discover all attributes and record them with their type
	attrMap := map[string]string{}
	if rows, err = tx.Query(
		stmt.NamespaceAttributeDiscover,
		q.Runtime.Namespace,
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

	// for all properties, any currently valid value must
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
			// the ID of the server being edited
			rteID,
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
				rteID,
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

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Runtime = append(mr.Runtime, q.Runtime)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
