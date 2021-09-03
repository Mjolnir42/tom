/*-
 * Copyright (c) 2016-2019, 1&1 Internet SE
 * Copyright (c) 2016-2019, Jörg Pernfuß <joerg.pernfuss@code.jpe@gmail.com>
 * Copyright (c) 2021, Jörg Pernfuß <joerg.pernfuss@ionos.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cmpl

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func GenericMulti(c *cli.Context, singlewords, multiwords []string) {
	keywords := append(singlewords, multiwords...)

	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range singlewords {
		if !match[t] {
			fmt.Println(t)
		}
	}
	for _, t := range multiwords {
		fmt.Println(t)
	}
}

// GenericMultiTriple tab-completes arguments in triples of
// (keyword, key, value)
func GenericMultiTriple(c *cli.Context, singlewords, multiwords []string) {
	keywords := append(singlewords, multiwords...)

	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 2
		match[t] = true
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range singlewords {
		if !match[t] {
			fmt.Println(t)
		}
	}
	for _, t := range multiwords {
		fmt.Println(t)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
