/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

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
func (h *NamespaceWriteHandler) txPropUpdate(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
	txTime *time.Time,
	prop proto.PropertyDetail,
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
	case `always`:
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
	case `forever`:
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

	// select standard|unique
	var stmtSelect, stmtClamp, stmtAdd *sql.Stmt
	switch attributeType {
	case `standard`:
		stmtSelect = tx.Stmt(h.stmtTxStdPropSelect)
		stmtClamp = tx.Stmt(h.stmtTxStdPropClamp)
		stmtAdd = tx.Stmt(h.stmtTxStdPropAdd)
	case `unique`:
		stmtSelect = tx.Stmt(h.stmtTxUniqPropSelect)
		stmtClamp = tx.Stmt(h.stmtTxUniqPropClamp)
		stmtAdd = tx.Stmt(h.stmtTxUniqPropAdd)
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
	); !ok {
		return false
	}
	// txPropClamp found an existing record that has the same value as well as
	// the same validUntil timestamp -> nothing to do further
	if done {
		return true
	}

	// txPropClamp either found no pre-existing value, or the existing value
	// was successfully invalidated. Set the new value as requested.
	if ok = h.txPropSetValue(
		q,
		mr,
		stmtAdd,
		prop,
		reqValidSince,
		reqValidUntil,
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
func (h *NamespaceWriteHandler) txPropClamp(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
	txTime *time.Time,
	stmtSelect *sql.Stmt,
	stmtClamp *sql.Stmt,
	prop proto.PropertyDetail,
	reqValidSince time.Time,
	reqValidUntil time.Time,
) (ok bool, done bool) {
	var (
		value                 string
		validFrom, validUntil time.Time
		res                   sql.Result
		err                   error
	)

	// query if the attribute has a value that is valid at transaction time
	if err := stmtSelect.QueryRow(
		q.Namespace.Name,
		prop.Attribute,
		txTime.Format(msg.RFC3339Milli),
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
		q.Namespace.Name,
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
func (h *NamespaceWriteHandler) txAttrQueryType(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
	prop proto.PropertyDetail,
) (typ string, ok bool) {
	if err := tx.Stmt(h.stmtAttQueryType).QueryRow(
		q.Namespace.Name,
		prop.Attribute,
	).Scan(
		&typ,
	); err == sql.ErrNoRows {
		mr.ServerError(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	ok = true
	return
}

// txPropSetValue realizes a property by setting a new attribute value
func (h *NamespaceWriteHandler) txPropSetValue(
	q *msg.Request,
	mr *msg.Result,
	stmt *sql.Stmt,
	prop proto.PropertyDetail,
	reqValidSince, reqValidUntil time.Time,
) bool {
	var res sql.Result
	var err error

	if res, err = stmt.Exec(
		q.Namespace.Name,
		prop.Attribute,
		prop.Value,
		reqValidSince,
		reqValidUntil,
		q.UserIDLib,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err)
		return false
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
