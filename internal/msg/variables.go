/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"errors"
	"time"
)

var (
	// this will be used as mapping for the PostgreSQL time value
	// -infinity. Dates earlier than this will be truncated to NegTimeInf
	//
	// Earliest value supported by postgreSQL: 4713 BC
	//
	// RFC3339: -4096-01-01T00:00:00Z
	NegTimeInf = time.Date(-4096, time.January, 1, 0, 0, 0, 0, time.UTC)
	// this will be used as mapping for the PostgreSQL time value
	// +infinity. Dates after this will be truncated to PosTimeInf.
	//
	// Latest value supported by postgreSQL: 294276 AD
	//
	// RFC: 293888-01-01T00:00:00Z
	PosTimeInf = time.Date(293888, time.January, 1, 0, 0, 0, 0, time.UTC)
	//
	ErrInvalidValidity = errors.New(`msg: invalid lower or upper validity boundary`)
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
