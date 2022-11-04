/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cred // import "github.com/mjolnir42/tom/internal/cred"

import (
	"math/bits"
	"strings"
)

func mask(pass []byte) []byte {
	pass = []byte(strings.TrimSpace(string(pass)))

	for i := range pass {
		// ROT13 ;)
		pass[i] = pass[i] | 0b00001101
		// mask used password vs disk stored password
		pass[i] = pass[i] ^ byte(0b01001101)

		pass[i] = bits.Reverse8(pass[i])
	}
	return pass
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
