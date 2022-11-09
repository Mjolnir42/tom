/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"fmt"
	"time"
)

func ResolveValidSince(s string, t, tx *time.Time) (err error) {
	switch s {
	case `always`:
		*t = NegTimeInf
	case `forever`, `perpetual`:
		err = fmt.Errorf("Invalid keyword for ValidSince: %s", s)
	case ``, `now`:
		*t = *tx
	default:
		*t, err = time.Parse(
			RFC3339Milli,
			s,
		)
	}
	return
}

func ResolvePValidSince(s string, t, tx *time.Time) (err error) {
	switch s {
	case `always`, `perpetual`:
		*t = NegTimeInf
	case `forever`:
		err = fmt.Errorf("Invalid keyword for ValidSince: %s", s)
	case ``, `now`:
		*t = *tx
	default:
		*t, err = time.Parse(
			RFC3339Milli,
			s,
		)
	}
	return
}

func ResolveValidUntil(s string, t, tx *time.Time) (err error) {
	switch s {
	case `always`, `perpetual`:
		err = fmt.Errorf("Invalid keyword for ValidSince: %s", s)
	case `forever`:
		*t = PosTimeInf
	case ``:
		*t = PosTimeInf
	case `now`:
		*t = *tx
	default:
		*t, err = time.Parse(
			RFC3339Milli,
			s,
		)
	}
	return
}

func ResolvePValidUntil(s string, t, tx *time.Time) (err error) {
	switch s {
	case `always`:
		err = fmt.Errorf("Invalid keyword for ValidSince: %s", s)
	case `forever`, `perpetual`:
		*t = PosTimeInf
	case ``:
		*t = PosTimeInf
	case `now`:
		*t = *tx
	default:
		*t, err = time.Parse(
			RFC3339Milli,
			s,
		)
	}
	return
}

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

func ParseValidUntil(s string, txNow *time.Time) (until *time.Time, err error) {
	// no consistent transaction timestamp for now was supplied
	if txNow == nil {
		*(txNow) = time.Now().UTC()
	}

	switch s {
	case `always`:
		err = ErrInvalidValidity
	case ``, `forever`, `perpetual`:
		*(until) = PosTimeInf
	case `now`:
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
