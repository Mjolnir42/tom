/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"errors"
)

var (
	ErrEmptyTomID   = errors.New(`proto: .TomID struct field is empty`)
	ErrInvalidTomID = errors.New(`proto: .TomID struct field contents are invalid`)
	ErrParsingTomID = errors.New(`Argument can not be parsed as a valid TomID`)
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
