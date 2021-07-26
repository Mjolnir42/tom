/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"time"

	//	"github.com/mjolnir42/lhm"
	//	"github.com/mjolnir42/tom/internal/handler"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/pkg/proto"
	//	"github.com/mjolnir42/tom/internal/stmt"
)

// process is the request dispatcher
func (h *NamespaceWriteHandler) process(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case msg.ActionAdd:
		h.add(q, &result)
	case msg.ActionRemove:
		h.remove(q, &result)
	case msg.ActionAttrAdd:
		h.attributeAdd(q, &result)
	case msg.ActionAttrRemove:
		h.attributeRemove(q, &result)
	case msg.ActionPropSet:
		h.propertySet(q, &result)
	case msg.ActionPropUpdate:
		h.propertyUpdate(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add creates a new namespace
func (h *NamespaceWriteHandler) add(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
		tx  *sql.Tx
	)

	//
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	//
	if res, err = tx.Stmt(h.stmtAdd).Exec(
		q.Namespace.Property[`dict_name`].Value,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}
	if !mr.CheckRowsAffected(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	//
	for _, property := range []string{`dict_type`, `dict_lookup`, `dict_uri`, `dict_ntt_list`} {
		if _, ok := q.Namespace.Property[property]; ok {
			if res, err = tx.Stmt(h.stmtConfig).Exec(
				q.Namespace.Property[`dict_name`].Value,
				property,
				q.Namespace.Property[property].Value,
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

	//
	for _, attribute := range q.Namespace.Attributes {
		if attribute.Unique {
			res, err = tx.Stmt(h.stmtAttUnqAdd).Exec(
				q.Namespace.Property[`dict_name`].Value,
				attribute.Key,
			)
		} else {
			res, err = tx.Stmt(h.stmtAttStdAdd).Exec(
				q.Namespace.Property[`dict_name`].Value,
				attribute.Key,
			)
		}
		if err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		if !mr.CheckRowsAffected(res.RowsAffected()) {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// remove deletes a specific namespace
func (h *NamespaceWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// attributeAdd ...
func (h *NamespaceWriteHandler) attributeAdd(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
		tx  *sql.Tx
	)

	//
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	//
	for _, attribute := range q.Namespace.Attributes {
		if attribute.Unique {
			res, err = tx.Stmt(h.stmtAttUnqAdd).Exec(
				q.Namespace.Name,
				attribute.Key,
			)
		} else {
			res, err = tx.Stmt(h.stmtAttStdAdd).Exec(
				q.Namespace.Name,
				attribute.Key,
			)
		}
		if err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		if !mr.CheckRowsAffected(res.RowsAffected()) {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// attributeRemove ...
func (h *NamespaceWriteHandler) attributeRemove(q *msg.Request, mr *msg.Result) {
	mr.NotImplemented() // TODO
}

// propertySet ...
func (h *NamespaceWriteHandler) propertySet(q *msg.Request, mr *msg.Result) {
	// tx.Begin()
	// forall Property in Request
	// 		NamespaceTxStdPropertySelect -> already has value?
	// 		yes: NamespaceTxStdPropertyClamp
	// 		NamespaceTxStdPropertyAdd
	// forall Property currently set but not in Request
	//		NamespaceTxStdPropertyClamp
	// tx.Commit()
	var (
		err    error
		tx     *sql.Tx
		txTime time.Time
		ok     bool
		rows   *sql.Rows
	)
	txTime = time.Now().UTC()
	attrMap := map[string]string{}

	// tx.Begin()
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// discover all attributes
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

	// forall Property in Request
	for key := range q.Namespace.Property {
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Namespace.Property[key],
		); !ok {
			tx.Rollback()
			return
		}
		// remove property contained in request from map of attributes
		delete(attrMap, q.Namespace.Property[key].Attribute)
	}

	// forall Property not in request
	for key := range attrMap {
		var stmtSelect, stmtClamp *sql.Stmt

		switch attrMap[key] {
		case `standard`:
			stmtSelect = tx.Stmt(h.stmtTxStdPropSelect)
			stmtClamp = tx.Stmt(h.stmtTxStdPropClamp)
		case `unique`:
			stmtSelect = tx.Stmt(h.stmtTxUniqPropSelect)
			stmtClamp = tx.Stmt(h.stmtTxUniqPropClamp)
		default:
			mr.ServerError()
			tx.Rollback()
			return
		}

		if ok, _ = h.txPropClamp(
			q, mr, tx, &txTime, stmtSelect, stmtClamp,
			proto.PropertyDetail{
				Attribute: key,
				Value:     txTime.Format(msg.RFC3339Milli) + key + `_clamp`,
			},
			txTime,
			txTime,
		); !ok {
			tx.Rollback()
			return
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

// propertyUpdate ...
func (h *NamespaceWriteHandler) propertyUpdate(q *msg.Request, mr *msg.Result) {
	var (
		err    error
		tx     *sql.Tx
		txTime time.Time
		ok     bool
	)
	txTime = time.Now().UTC()

	// tx.Begin()
	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	// forall Property in Request
	for key := range q.Namespace.Property {
		if ok = h.txPropUpdate(
			q, mr, tx, &txTime, q.Namespace.Property[key],
		); !ok {
			tx.Rollback()
			return
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

// txPropUpdate ...
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
	if done {
		return true
	}

	// 		NamespaceTxStdPropertyAdd
	if ok = h.txPropSetValue(
		q,
		mr,
		tx,
		stmtAdd,
		prop,
		reqValidSince,
		reqValidUntil,
	); !ok {
		return false
	}

	return true
}

// txAttrQueryType ...
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

// txPropSetValue ...
func (h *NamespaceWriteHandler) txPropSetValue(
	q *msg.Request,
	mr *msg.Result,
	tx *sql.Tx,
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
		reqValidSince.Format(msg.RFC3339Milli),
		reqValidUntil.Format(msg.RFC3339Milli),
	); err != nil {
		mr.ServerError(err)
		return false
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return false
	}
	return true
}

// txPropClamp ...
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
	//	NamespaceTxStdPropertySelect -> already has value?
	var (
		value                 string
		validFrom, validUntil time.Time
		res                   sql.Result
		err                   error
	)

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

	// 		yes:NamespaceTxStdPropertyClamp
	// moving lower validity bound to earlier is unsupported
	if reqValidSince.Before(validFrom) {
		mr.BadRequest()
		return false, false
	}
	if reqValidUntil.Equal(validUntil) && value == prop.Value {
		// nothing to do
		return true, true
	}

	if res, err = stmtClamp.Exec(
		q.Namespace.Name,
		prop.Attribute,
		value,
		validFrom.Format(msg.RFC3339Milli),
		validUntil.Format(msg.RFC3339Milli),
		txTime.Format(msg.RFC3339Milli),
		reqValidSince.Format(msg.RFC3339Milli),
	); err != nil {
		mr.ServerError(err)
		return false, false
	}
	if !mr.AssertOneRowAffected(res.RowsAffected()) {
		return false, false
	}
	return true, false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
