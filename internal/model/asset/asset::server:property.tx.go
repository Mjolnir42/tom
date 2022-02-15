/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package asset // import "github.com/mjolnir42/tom/internal/model/asset"

import (
	"database/sql"
	"time"

	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
)

// this file contains transaction helper functions for property requests

// txPropUpdate ensures that the specified property is set, ie. that the
// attribute indicated in prop has the value specified in prop during the
// specified period of validity.
// If prop.ValidSince is not specifed, the new value starts its validity at
// transaction time txTime. If prop.ValidUntil is not specified, then the
// new value will be set to valid forever.
func (h *ServerWriteHandler) txPropUpdate(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
	txTime *time.Time,
	prop proto.PropertyDetail,
	serverID string,
) bool {
	var (
		attributeType string
		ok, done      bool
		err           error
	)

	// NamespaceAttributeQueryType -> standard|unique
	if attributeType, ok = h.txAttrQueryType(
		q, mr, tx, prop,
	); !ok {
		return ok
	}

	var reqValidSince, reqValidUntil time.Time
	switch prop.ValidSince {
	case `always`, `perpetual`:
		reqValidSince = msg.NegTimeInf
	case `forever`:
		mr.BadRequest()
		return false
	case ``:
		reqValidSince = *txTime
	default:
		if reqValidSince, err = time.Parse(
			msg.RFC3339Milli,
			prop.ValidSince,
		); err != nil {
			mr.BadRequest(err)
			return false
		}
	}

	switch prop.ValidUntil {
	case `always`:
		mr.BadRequest()
		return false
	case `forever`, `perpetual`:
		reqValidUntil = msg.PosTimeInf
	case ``:
		reqValidUntil = msg.PosTimeInf
	default:
		if reqValidUntil, err = time.Parse(
			msg.RFC3339Milli,
			prop.ValidUntil,
		); err != nil {
			mr.BadRequest(err)
			return false
		}
	}

	// check the use of the perpetual keyword
	if prop.ValidSince == `perpetual` || prop.ValidUntil == `perpetual` {
		// both Since and Until must be set to perpetual
		if prop.ValidSince != prop.ValidUntil {
			mr.BadRequest()
			return false
		}
		// only the `type` property is perpetual
		if prop.Attribute != `type` {
			mr.BadRequest()
			return false
		}
	}

	// select standard|unique
	var stmtSelect, stmtClamp, stmtAdd, stmtClean *sql.Stmt
	switch attributeType {
	case proto.AttributeStandard:
		stmtSelect = tx.Stmt(h.stmtTxStdPropSelect)
		stmtClamp = tx.Stmt(h.stmtTxStdPropClamp)
		stmtAdd = tx.Stmt(h.stmtTxStdPropAdd)
		stmtClean = tx.Stmt(h.stmtTxStdPropClean)
	case proto.AttributeUnique:
		stmtSelect = tx.Stmt(h.stmtTxUniqPropSelect)
		stmtClamp = tx.Stmt(h.stmtTxUniqPropClamp)
		stmtAdd = tx.Stmt(h.stmtTxUniqPropAdd)
		stmtClean = tx.Stmt(h.stmtTxUniqPropClean)
	default:
		mr.ServerError()
		return false
	}

	// check if there is an existing, at transaction time valid value
	// associated with the attribute, and clamp its validity down so a new
	// value can be set.
	if ok, done = h.txPropClamp(
		q,
		mr,
		tx,
		txTime,
		stmtSelect,
		stmtClamp,
		prop,
		reqValidSince,
		reqValidUntil,
		serverID,
	); !ok {
		return false
	}
	// txPropClamp found an existing record that has the same value as well as
	// the same validUntil timestamp -> nothing to do further
	if done {
		return true
	}

	// txPropClamp either found no pre-existing value, or the existing value
	// was successfully invalidated. Clean future records that have never
	// been valid yet, ie. their validSince is after now()
	if ok = h.txPropClean(
		q,
		mr,
		stmtClean,
		prop,
		txTime,
		serverID,
	); !ok {
		return false
	}

	// set the new value as requested.
	if ok = h.txPropSetValue(
		q,
		mr,
		stmtAdd,
		prop,
		reqValidSince,
		reqValidUntil,
		serverID,
	); !ok {
		return false
	}

	return true
}

// txPropClamp examines if there is a pre-existing valid value for an
// attribute. If no pre-existing value exists, or the validUntil of the
// existing value was successfully clamped down to the validSince of
// the new value, the function returns ok.
// When txPropClamp returns !ok, the error value of mr is always set.
// If the new value matches the old value and the new validUntil matches
// the old validUntil, then the record is left and place and the function
// indicates via the done indicator that no new record needs to be inserted.
func (h *ServerWriteHandler) txPropClamp(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
	txTime *time.Time,
	stmtSelect *sql.Stmt,
	stmtClamp *sql.Stmt,
	prop proto.PropertyDetail,
	reqValidSince time.Time,
	reqValidUntil time.Time,
	serverID string,
) (ok bool, done bool) {
	var (
		value                 string
		validFrom, validUntil time.Time
		res                   sql.Result
		err                   error
	)

	// query if the attribute has a value that is valid at transaction time
	if err := stmtSelect.QueryRow(
		q.Server.Namespace,
		prop.Attribute,
		txTime.Format(msg.RFC3339Milli),
		serverID,
	).Scan(
		&value,
		&validFrom,
		&validUntil,
	); err == sql.ErrNoRows {
		// not having a value is not an error
		return true, false
	} else if err != nil {
		mr.ServerError(err)
		return false, false
	}

	// the lower validity boundary of the existing record is not updated
	// by this request, so if the requested lower validity boundary for the
	// new record is before the lower validity boundary of the old record,
	// then the update can not be performed.
	// This would mean that the higher validity boundary of the existing
	// record must be set to before the lower boundary.
	if reqValidSince.Before(validFrom) {
		mr.BadRequest()
		return false, false
	}
	if reqValidUntil.Equal(validUntil) && value == prop.Value {
		// nothing to do, the new value to be set matches the already existing
		// record
		return true, true
	}

	if res, err = stmtClamp.Exec(
		// name of the namespace
		q.Server.Namespace,
		// name of the attribute
		prop.Attribute,
		// current value, required for row confirmation
		value,
		// current validFrom, required for row confirmation
		validFrom,
		// current validUntil, required for row confirmation
		validUntil,
		// updated row must be valid at txTime
		txTime,
		// update upper validity of the existing record to lower validity
		// of the new record
		reqValidSince,
		// ID of the server environment
		serverID,
	); err != nil {
		mr.ServerError(err)
		return false, false
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		// error set by mr.AssertOneRowAffected
		return false, false
	}
	return true, false
}

// txAttrQueryType returns the type for a single attribute
func (h *ServerWriteHandler) txAttrQueryType(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
	prop proto.PropertyDetail,
) (typ string, ok bool) {
	if err := tx.Stmt(h.stmtAttQueryType).QueryRow(
		q.Server.Namespace,
		prop.Attribute,
	).Scan(
		&typ,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	ok = true
	return
}

// txPropSetValue realizes a property by setting a new attribute value
func (h *ServerWriteHandler) txPropSetValue(
	q *msg.Request,
	mr *msg.Result,
	stmt *sql.Stmt,
	prop proto.PropertyDetail,
	reqValidSince, reqValidUntil time.Time,
	serverID string,
) bool {
	var res sql.Result
	var err error

	if res, err = stmt.Exec(
		q.Server.Namespace,
		prop.Attribute,
		prop.Value,
		reqValidSince,
		reqValidUntil,
		q.UserIDLib,
		q.AuthUser,
		serverID,
	); err != nil {
		mr.ServerError(err)
		return false
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return false
	}
	return true
}

// txPropClean deletes all records with a starting validity in the
// future. This is restricted to current time, to ensure that no records
// that were valid in the past can be deleted.
func (h *ServerWriteHandler) txPropClean(
	q *msg.Request,
	mr *msg.Result,
	stmt *sql.Stmt,
	prop proto.PropertyDetail,
	txTime *time.Time,
	serverID string,
) bool {
	var err error

	if _, err = stmt.Exec(
		q.Server.Namespace,
		prop.Attribute,
		serverID,
		*txTime,
	); err != nil {
		mr.ServerError(err)
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
