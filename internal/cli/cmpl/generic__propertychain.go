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

func GenericMultiWithProperty(c *cli.Context, singlewords, multiwords []string) {
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

	for i, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		if t == `property` {
			propertySubChain(c.Args().Tail()[i:], singlewords, multiwords)
			return
		}
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

// GenericPropertyChain tab-completes with property chains
func GenericPropertyChain(c *cli.Context, singlewords, multiwords []string) {
	keywords := append(singlewords, multiwords...)
	chainwords := []string{`value`, `since`, `until`}

	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	inChain := false

	skip := 0
	match := make(map[string]bool)
	chainMatch := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		if t == `property` {
			inChain = true
			chainMatch = map[string]bool{}
			chainMatch[t] = true
		} else if inChain {
			switch t {
			case `value`:
				chainMatch[t] = true
			case `since`:
				chainMatch[t] = true
			case `until`:
				chainMatch[t] = true
			default:
				inChain = false
			}
		}
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	// within the property chain, only complete chain arguments
	if inChain {
		for _, t := range chainwords {
			if !chainMatch[t] {
				fmt.Println(t)
			}
		}
		if chainMatch[`value`] {
			fmt.Println(`property`)
		}
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

func propertySubChain(args []string, singlewords, multiwords []string) {
	chainwords := []string{`value`, `since`, `until`}

	inChain := false

	skip := 0
	match := make(map[string]bool)
	chainMatch := make(map[string]bool)

	for _, t := range args {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		if t == `property` {
			inChain = true
			chainMatch = map[string]bool{}
			chainMatch[t] = true
		} else if inChain {
			switch t {
			case `value`:
				chainMatch[t] = true
			case `since`:
				chainMatch[t] = true
			case `until`:
				chainMatch[t] = true
			default:
				inChain = false
			}
		}
		continue
	}
	// do not complete in positions where arguments are expected
	if skip > 0 {
		return
	}
	// within the property chain, only complete chain arguments
	if inChain {
		for _, t := range chainwords {
			if !chainMatch[t] {
				fmt.Println(t)
			}
		}
		if chainMatch[`value`] {
			fmt.Println(`property`)
		}
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
