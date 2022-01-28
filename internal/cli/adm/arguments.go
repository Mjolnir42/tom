/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 * Copyright (c) 2021, Jörg Pernfuß <joerg.pernfuss@ionos.com>
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"
	"strings"

	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

// ParseVariadicArguments parses split up argument lists of
// keyword/value pairs were keywords can be specified multiple
// times, some keywords are required and some only allowed once.
// Sequence of multiple keywords are detected and lead to abort
//
//	multKeys => [ "port", "transport" ]
//	uniqKeys => [ "team" ]
//	reqKeys  => [ "team" ]
//	args     => [ "port", "53", "transport", "tcp", "transport",
//	              "udp", "team", "GenericOps" ]
//
//	result => result["team"] = [ "GenericOps" ]
//	          result["port"] = [ "53" ]
//	          result["transport"] = [ "tcp", "udp" ]
func ParseVariadicArguments(
	result map[string][]string, // provided result map
	multKeys []string, // keys that may appear multiple times
	uniqKeys []string, // keys that are allowed at most once
	reqKeys []string, // keys that are required at least one
	args []string, // arguments to parse
) error {
	// used to hold found errors, so if three keywords are missing they can
	// all be mentioned in one call
	errors := []string{}

	// merge key slices
	keys := append(multKeys, uniqKeys...)

	// helper to skip over next value in args slice
	skip := false

	for pos, val := range args {
		// skip current arg if last argument was a keyword
		if skip {
			skip = false
			continue
		}

		if sliceContainsString(val, keys) {
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				errors = append(errors,
					`Syntax error, incomplete key/value specification (too few items left to parse)`,
				)
				goto abort
			}
			// check for back-to-back keyswords
			if err := checkStringNotAKeyword(args[pos+1], keys); err != nil {
				errors = append(errors, err.Error())
				goto abort
			}

			// append value of current keyword into result map
			result[val] = append(result[val], strings.Trim(args[pos+1], `'"%`))
			skip = true
			continue
		}
		// keywords trigger continue before this
		// values after keywords are skip'ed
		// reaching this is an error
		errors = append(errors, fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		// ok is false if slice is nil
		if _, ok := result[key]; !ok {
			errors = append(errors, fmt.Sprintf("Syntax error, missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			errors = append(errors, fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

abort:
	if len(errors) > 0 {
		for key := range result {
			delete(result, key)
		}
		return fmt.Errorf(combineStrings(errors...))
	}

	return nil
}

// ParseVariadicTriples is a variant of ParseVariadicArguments where
// every keyword is followed by two values
func ParseVariadicTriples(
	result map[string][][2]string, // provided result map
	multKeys []string, // keys that may appear multiple times
	uniqKeys []string, // keys that are allowed at most once
	reqKeys []string, // keys that are required at least one
	args []string, // arguments to parse
) error {
	// used to hold found errors, so if three keywords are missing they can
	// all be mentioned in one call
	errors := []string{}

	// merge key slices
	keys := append(multKeys, uniqKeys...)

	// helper to skip over next value in args slice
	skip := false
	skipcount := 0

	for pos, val := range args {
		// skip current arg if last argument was a keyword
		if skip {
			skipcount--
			if skipcount == 0 {
				skip = false
			}
			continue
		}

		if sliceContainsString(val, keys) {
			// there must be at least two arguments left
			if len(args[pos+1:]) < 2 {
				errors = append(errors, `Syntax error, incomplete`+
					` key/value specification (too few items left`+
					` to parse)`,
				)
				goto abort
			}
			// check for back-to-back keyswords
			if err := checkStringNotAKeyword(args[pos+1],
				keys); err != nil {
				errors = append(errors, err.Error())
				goto abort
			}

			// append values of current keyword into result map
			result[val] = append(result[val],
				[2]string{
					strings.Trim(args[pos+1], `'"%`),
					args[pos+2],
				})
			skip = true
			skipcount = 2
			continue
		}
		// keywords trigger continue before this
		// values after keywords are skip'ed
		// reaching this is an error
		errors = append(errors, fmt.Sprintf("Syntax error, erroneus"+
			" argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		// ok is false if slice is nil
		if _, ok := result[key]; !ok {
			errors = append(errors, fmt.Sprintf("Syntax error,"+
				" missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			errors = append(errors, fmt.Sprintf("Syntax error,"+
				" keyword must only be provided once: %s", key))
		}
	}

abort:
	if len(errors) > 0 {
		for key := range result {
			delete(result, key)
		}
		return fmt.Errorf(combineStrings(errors...))
	}

	return nil
}

// ParsePropertyArguments ...
func ParsePropertyArguments(
	result map[string][]string,
	props *[]proto.PropertyDetail,
	args []string,
	variable []string,
	once []string,
	required []string,
) error {
	// used to hold found errors, so if three keywords are missing they can
	// all be mentioned in once call
	errors := []string{}

	// these are the keys inside a property chain
	propchain := []string{
		`value`,
		`since`,
		`until`,
	}

	// merge key slices
	keys := append(variable, once...)

	// helper to skip over next value in args slice
	skip := false
	skipcount := 0

argloop:
	for pos, val := range args {
		// skip current arg if last argument was a keyword
		if skip {
			skipcount--
			if skipcount == 0 {
				skip = false
			}
			continue argloop
		}

		if sliceContainsString(val, keys) {
			// current value is a keyword

			// there must be at least one argument entry left
			if len(args[pos+1:]) < 1 {
				errors = append(errors, `Syntax error, incomplete`+
					` key/value specification (too few items left`+
					` to parse)`,
				)
				goto abort
			}
			// check for back-to-back keywords
			if err := checkStringNotAKeyword(
				args[pos+1], keys,
			); err != nil {
				errors = append(errors, err.Error())
				goto abort
			}

			switch val {
			case `property`:
				// there must be at least three entries left
				if len(args[pos+1:]) < 3 {
					errors = append(errors, `Syntax error, incomplete`+
						` key/value specification (too few items left`+
						` to parse)`,
					)
					goto abort
				}
				// local copy of the prop keyword chain
				chain := make([]string, 0, len(propchain))
				chain = append(chain, propchain...)

				// check for back-to-back keywords inside the property chain
				if err := checkStringNotAKeyword(
					args[pos+1], propchain,
				); err != nil {
					errors = append(errors, err.Error())
					goto abort
				}

				// consume attribute name
				prop := proto.PropertyDetail{
					Attribute: strings.Trim(args[pos+1], `'"%`),
				}
				skip = true
				skipcount = 1

				// check next position, this must be a keyword from chain
				if sliceContainsString(args[pos+1+skipcount], chain) {
					// there must be at least one argument entry left
					if len(args[pos+1+skipcount:]) < 1 {
						errors = append(errors, `Syntax error, incomplete`+
							` key/value specification (too few items left`+
							` to parse)`,
						)
						goto abort
					}
					i := elemNumInSlice(chain, args[pos+1+skipcount])
					switch args[pos+1+skipcount] {
					case `value`:
						prop.Value = args[pos+1+skipcount+1]
					case `since`:
						prop.ValidSince = args[pos+1+skipcount+1]
					case `until`:
						prop.ValidUntil = args[pos+1+skipcount+1]
					}
					skipcount += 2
					chain = append(chain[:i], chain[i+1:]...)
				} else {
					// at least `property attr value val` is required
					errors = append(errors, `Syntax error, incomplete`+
						` specification of a property`,
					)
					goto abort
				}
				if len(args[pos+1+skipcount:]) == 0 {
					if prop.Value != `` {
						*props = append(*props, prop)
					}
					continue argloop
				}

				// check next position, this MAY be a keyword from chain
				if sliceContainsString(args[pos+1+skipcount], chain) {
					// there must be at least one argument entry left
					if len(args[pos+1+skipcount:]) < 1 {
						errors = append(errors, `Syntax error, incomplete`+
							` key/value specification (too few items left`+
							` to parse)`,
						)
						goto abort
					}
					i := elemNumInSlice(chain, args[pos+1+skipcount])
					switch args[pos+1+skipcount] {
					case `value`:
						prop.Value = args[pos+1+skipcount+1]
					case `since`:
						prop.ValidSince = args[pos+1+skipcount+1]
					case `until`:
						prop.ValidUntil = args[pos+1+skipcount+1]
					}
					skipcount += 2
					chain = append(chain[:i], chain[i+1:]...)
				} else if sliceContainsString(args[pos+1+skipcount], keys) {
					// this position resumes with regular keywords,
					// break property chain parsing if value was already found
					if prop.Value != `` {
						*props = append(*props, prop)
						continue argloop
					}
					errors = append(errors, `Syntax error, incomplete`+
						` specification of a property`,
					)
					goto abort
				} else {
					if prop.Value != `` {
						*props = append(*props, prop)
					}
					// this will trigger the error for not-a-keyword
					continue argloop
				}
				if len(args[pos+1+skipcount:]) == 0 {
					if prop.Value != `` {
						*props = append(*props, prop)
					}
					continue argloop
				}

				// check next position, this MAY be a keyword from chain
				if sliceContainsString(args[pos+1+skipcount], chain) {
					// there must be at least one argument entry left
					if len(args[pos+1+skipcount:]) < 1 {
						errors = append(errors, `Syntax error, incomplete`+
							` key/value specification (too few items left`+
							` to parse)`,
						)
						goto abort
					}
					i := elemNumInSlice(chain, args[pos+1+skipcount])
					switch args[pos+1+skipcount] {
					case `value`:
						prop.Value = args[pos+1+skipcount+1]
					case `since`:
						prop.ValidSince = args[pos+1+skipcount+1]
					case `until`:
						prop.ValidUntil = args[pos+1+skipcount+1]
					}
					skipcount += 2
					chain = append(chain[:i], chain[i+1:]...)
					if prop.Value != `` {
						*props = append(*props, prop)
					}
				} else if sliceContainsString(args[pos+1+skipcount], keys) {
					// this position resumes with regular keywords,
					// break property chain parsing if value was already found
					if prop.Value != `` {
						*props = append(*props, prop)
						continue argloop
					}
					errors = append(errors, `Syntax error, incomplete`+
						` specification of a property`,
					)
					goto abort
				} else {
					if prop.Value != `` {
						*props = append(*props, prop)
					}
					// this will trigger the error for not-a-keyword
					continue argloop
				}

				continue argloop
			default:
				// regular key/value keyword
				result[val] = append(result[val], args[pos+1])
				skip = true
				skipcount = 1
				continue argloop
			}
		}
		// this error is reached if argument was not skipped and not a
		// recognized keyword
		errors = append(errors, fmt.Sprintf("Syntax error, erroneus"+
			" argument: %s", val))
	} // argloop

	// check if all required keywords were collected
requiredKeywordsLoop:
	for _, key := range required {
		if _, ok := result[key]; !ok {
			if key == `property` && len(*props) > 0 {
				continue requiredKeywordsLoop
			}
			errors = append(errors, fmt.Sprintf("Syntax error,"+
				" missing keyword: %s", key))
		}
	}

	for _, key := range once {
		// check ok since once may still be optional
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			errors = append(errors, fmt.Sprintf("Syntax error,"+
				" keyword must only be provided once: %s", key))
		}
	}

abort:
	if len(errors) > 0 {
		return fmt.Errorf(combineStrings(errors...))
	}
	return nil
}

// VerifySingleArgument takes a context and verifies there is only one
// commandline argument
func VerifySingleArgument(c *cli.Context) error {
	a := c.Args()
	if !a.Present() {
		return fmt.Errorf(`Syntax error, command requires argument`)
	}

	if len(a.Tail()) != 0 {
		return fmt.Errorf(
			"Syntax error, too many arguments (expected: 1, received %d)",
			len(a.Tail())+1,
		)
	}
	return nil
}

// VerifyNoArgument takes a context and verifies there is no
// commandline argument
func VerifyNoArgument(c *cli.Context) error {
	a := c.Args()
	if a.Present() {
		return fmt.Errorf(`Syntax error, command takes no arguments`)
	}

	return nil
}

// AllArguments returns all arguments from the given cli.Context
func AllArguments(c *cli.Context) []string {
	sl := []string{c.Args().First()}
	sl = append(sl, c.Args().Tail()...)

	if c.Args().First() == `` && len(sl) == 1 {
		return []string{}
	}

	return sl
}

// sliceContainsString checks whether string s is in slice sl
func sliceContainsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

// checkStringNotAKeyword checks whether string s in not in slice keys
func checkStringNotAKeyword(s string, keys []string) error {
	if sliceContainsString(s, keys) {
		return fmt.Errorf("Syntax error, back-to-back keyword: %s", s)
	}
	return nil
}

// combineStrings takes an arbitrary number of strings and combines them
// into one, separated by space.
func combineStrings(s ...string) string {
	var out string
	spacer := ``
	for _, in := range s {
		// ensure a single trailing .
		out = fmt.Sprintf("%s%s", out+spacer, strings.TrimRight(in, `.`)+`.`)
		spacer = ` `
	}
	return out
}

// VariadicArguments calls ParseVariadicArguments with a predefined
// argument configuration suitable for a given command
func VariadicArguments(command string, c *cli.Context, result *map[string][]string) error {

	multipleAllowed, uniqueOptions, mandatoryOptions := ArgumentsForCommand(command)
	return ParseVariadicArguments(
		*result,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	)
}

func VariadicDirect(command string, c *cli.Context, result *map[string][]string) error {

	multipleAllowed, uniqueOptions, mandatoryOptions := ArgumentsForCommand(command)
	var args []string
	// for no arguments, c.Args().First() will return an empty string. A
	// slice of strings with an empty string inside is not the same as
	// an empty slice of strings
	if c.Args().First() != `` {
		args = append([]string{c.Args().First()}, c.Args().Tail()...)
	}
	return ParseVariadicArguments(
		*result,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		args,
	)
}

// VariadicTriples calls ParseVariadicTriples with a predefined
// argument configuration suitable for a given command
func VariadicTriples(command string, c *cli.Context, result *map[string][][2]string) error {

	multipleAllowed, uniqueOptions, mandatoryOptions := ArgumentsForCommand(command)
	return ParseVariadicTriples(
		*result,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	)
}

func elemNumInSlice(sl []string, el string) int {
	for i, v := range sl {
		if v == el {
			return i
		}
	}
	return -1
}

// ArgumentsForCommand contains the database of argument configurations
// for VariadicArguments
func ArgumentsForCommand(s string) (multipleAllowed, uniqueOptions, mandatoryOptions []string) {
	switch s {
	case proto.CmdNamespaceAdd:
		return []string{`std-attr`, `uniq-attr`}, []string{`type`, `lookup-uri`, `lookup-key`}, []string{`type`}
	case proto.CmdNamespaceAttrAdd, proto.CmdNamespaceAttrRemove:
		return []string{`std-attr`, `uniq-attr`}, []string{}, []string{}
	case proto.CmdNamespacePropSet, proto.CmdNamespacePropUpdate, proto.CmdNamespacePropRemove:
		return []string{`property`}, []string{}, []string{`property`}
	case proto.CmdContainerAdd:
		fallthrough
	case proto.CmdServerAdd:
		fallthrough
	case proto.CmdOrchestrationAdd:
		fallthrough
	case proto.CmdRuntimeAdd:
		return []string{`property`}, []string{`namespace`, `type`, `since`, `until`}, []string{`namespace`, `type`}
	case proto.CmdContainerList, proto.CmdContainerShow, proto.CmdContainerRemove:
		fallthrough
	case proto.CmdServerList, proto.CmdServerShow, proto.CmdServerRemove:
		fallthrough
	case proto.CmdRuntimeList, proto.CmdRuntimeShow, proto.CmdRuntimeRemove:
		return []string{}, []string{`namespace`}, []string{}
	case proto.CmdContainerPropSet, proto.CmdContainerPropUpdate, proto.CmdContainerPropRemove:
		fallthrough
	case proto.CmdServerPropSet, proto.CmdServerPropUpdate, proto.CmdServerPropRemove:
		fallthrough
	case proto.CmdRuntimePropSet, proto.CmdRuntimePropUpdate, proto.CmdRuntimePropRemove:
		return []string{`property`}, []string{`namespace`}, []string{`property`, `namespace`}
	case proto.CmdContainerLink:
		fallthrough
	case proto.CmdServerLink:
		fallthrough
	case proto.CmdRuntimeLink:
		return []string{`is-equal`}, []string{}, []string{`is-equal`}
	case proto.CmdServerStack:
		return []string{}, []string{`namespace`, `since`, `until`, `provided-by`}, []string{`provided-by`}
	case proto.CmdRuntimeStack:
		return []string{}, []string{`namespace`, `runs-on`}, []string{`runs-on`}
	default:
		return []string{}, []string{}, []string{}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
