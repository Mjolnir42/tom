/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"time"
)

//
func ParseValidSince(s string, txNow *time.Time) (since *time.Time, err error) {
	// no consistent transaction timestamp for now was supplied
	if txNow == nil {
		*(txNow) = time.Now().UTC()
	}

	switch s {
	case `always`, `perpetual`:
		*(since) = NegTimeInf
	case `forever`:
		err = ErrInvalidValidity
	case ``, `now`:
		since = txNow
	default:
		*(since), err = time.Parse(
			RFC3339Milli,
			s,
		)
	}
	return since, err
}

//
func ParseValidUntil(s string, txNow *time.Time) (until *time.Time, err error) {
	// no consistent transaction timestamp for now was supplied
	if txNow == nil {
		*(txNow) = time.Now().UTC()
	}

	switch s {
	case `always`:
		err = ErrInvalidValidity
	case `forever`, `perpetual`:
		*(until) = PosTimeInf
	case ``, `now`:
		until = txNow
	default:
		*(until), err = time.Parse(
			RFC3339Milli,
			s,
		)
	}
	return until, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
