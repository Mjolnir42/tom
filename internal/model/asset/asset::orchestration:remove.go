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
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdOrchestrationRemove)

	registry = append(registry, function{
		cmd:    proto.CmdOrchestrationRemove,
		handle: orchestrationRemove,
	})
}

func orchestrationRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.OrchestrationRemove)
}

func exportOrchestrationRemove(result *proto.Result, r *msg.Result) {
	result.Orchestration = &[]proto.Orchestration{}
	*result.Orchestration = append(*result.Orchestration, r.Orchestration...)
}

// OrchestrationRemove function
func (m *Model) OrchestrationRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionOrchestration
	request.Action = proto.ActionRemove
	request.Orchestration = *(proto.NewOrchestration())
	request.Orchestration.TomID = params.ByName(`tomID`)

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
	m.x.Send(&w, &result, exportOrchestrationRemove)
}

// remove invalidates an existing orchestration
func (h *OrchestrationWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	var (
		txTime                                           time.Time
		oreID, dictionaryID, createdAt, createdBy        string
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

	// discover oreID at the start of the transaction
	if err = tx.Stmt(
		h.stmtTxShow,
	).QueryRow(
		q.Orchestration.Namespace,
		q.Orchestration.Name,
		txTime,
	).Scan(
		&oreID,
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

	// unstack the orchestration itself
	txStackClamp := tx.Stmt(h.stmtTxStackClampAll)
	if _, err = txStackClamp.Exec(
		txTime,
		oreID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// unstack all runtime children currently on the orchestration
	if _, err = tx.Stmt(
		h.stmtTxChild,
	).Exec(
		txTime,
		oreID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if _, err = tx.Stmt(
		h.stmtTxChildClean,
	).Exec(
		txTime,
		oreID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// discover all attributes and record them with their type
	attrMap := map[string]string{}
	if rows, err = tx.Stmt(
		h.stmtAttDiscover,
	).Query(
		q.Orchestration.Namespace,
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
			oreID,
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
				oreID,
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
	mr.Orchestration = append(mr.Orchestration, q.Orchestration)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
