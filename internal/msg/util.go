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
			*t = t.UTC()
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
			*t = t.UTC()
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
			*t = t.UTC()
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
			*t = t.UTC()
	}
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
