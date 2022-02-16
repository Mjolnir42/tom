/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// OnlyUnreserved checks that s only contains characters from proto.CharUnreserved
func OnlyUnreserved(s string) error {
	for _, b := range []byte(s) {
		if !strings.Contains(CharUnreserved, string(b)) {
			return fmt.Errorf(
				"String <%s> contains illegal character <%s>",
				s,
				string(b),
			)
		}
	}
	return nil
}

// ValidNamespace checks that s only contains characters and structure suitable
// for a namespace name
func ValidNamespace(s string) error {
	for _, b := range []byte(s) {
		if !strings.Contains(CharNamespace, string(b)) {
			return fmt.Errorf(
				"String <%s> contains illegal character <%s>",
				s,
				string(b),
			)
		}
	}
	switch strings.Count(s, `~`) {
	case 0:
	case 1:
		switch {
		case strings.HasPrefix(s, `tool~`):
		case strings.HasPrefix(s, `team~`):
		default:
			return fmt.Errorf(
				"Unknown namespace prefix: %s",
				strings.Split(s, `~`)[1],
			)
		}
	default:
		return fmt.Errorf(
			"Character <~> matches %d times but is only allowed once as prefix separator",
			strings.Count(s, `~`),
		)
	}
	return nil
}

// CheckPropertyConstraints checks that if the property attribute end
// with one of the suffixes _list, _json or _xml, then the value must
// decode appropriately.
func CheckPropertyConstraints(prop *PropertyDetail) error {
	if strings.HasSuffix(prop.Attribute, `_list`) {
		jl := &[]string{}
		if err := json.NewDecoder(bytes.NewBufferString(prop.Value)).Decode(jl); err != nil {
			return fmt.Errorf("Property %s is not a valid list: %s",
				prop.Attribute,
				err.Error(),
			)
		}
		return nil
	}
	if strings.HasSuffix(prop.Attribute, `_json`) {
		var j interface{}
		if err := json.NewDecoder(bytes.NewBufferString(prop.Value)).Decode(&j); err != nil {
			return fmt.Errorf("Property %s is not valid JSON: %s",
				prop.Attribute,
				err.Error(),
			)
		}
		return nil
	}
	if strings.HasSuffix(prop.Attribute, `_xml`) {
		var x interface{}
		if err := xml.Unmarshal([]byte(prop.Value), &x); err != nil {
			return fmt.Errorf("Property %s is not supported XML: %s",
				prop.Attribute,
				err.Error(),
			)
		}
		return nil
	}
	return nil
}

// AssertCommandIsDefined checks that the package variable Commands
// has an entry for command c
func AssertCommandIsDefined(c string) {
	if _, ok := Commands[c]; !ok {
		panic(c)
	}
}

// IsTomID returns true if s is a syntactically valid TomID
func IsTomID(s string) bool {
	return isTomIDFormatDNS(s) || isTomIDFormatURI(s)
}

func IsWildcardTomID(s string) bool {
	return isTomIDWildcardFormat(s)
}

// ParseTomID parses the TomID s and returns the entity type as a string as
// well as and Entity object
func ParseTomID(s string) (error, string, Entity) {
	var name, namespace, entity string
	switch {
	case s == ``:
		return ErrEmptyTomID, ``, nil
	case isTomIDFormatDNS(s):
		name, namespace, entity = parseTomIDFormatDNS(s)
	case isTomIDFormatURI(s):
		name, namespace, entity = parseTomIDFormatURI(s)
	default:
		return ErrInvalidTomID, ``, nil
	}

	switch entity {
	case EntityContainer:
		return nil, entity, (&Container{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntityNamespace:
		return nil, entity, (&Namespace{
			Name: name,
		}).SetTomID()
	case EntityOrchestration:
		return nil, entity, (&Orchestration{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntityRuntime:
		return nil, entity, (&Runtime{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntityServer:
		return nil, entity, (&Server{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntitySocket:
		return nil, entity, (&Socket{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	default:
		return ErrInvalidTomID, ``, nil
	}
}

func isTomIDFormatDNS(s string) bool {
	re := regexp.MustCompile(fmt.Sprintf("%s|%s|%s", tomIDFormatDNS, tomIDShortDNS, tomIDNamespDNS))
	return re.MatchString(s)
}

func isTomIDFormatURI(s string) bool {
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", tomIDFormatURI, tomIDNamespURI))
	return re.MatchString(s)
}

func isTomIDWildcardFormat(s string) bool {
	re := regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s",
		tomIDQueryNTT,
		tomIDSQueryNTT,
		tomIDQueryNsNTT,
		tomIDSQueryNsNTT,
	))
	return re.MatchString(s)
}

func parseTomIDFormatDNS(s string) (name, namespace, entity string) {
	reCommon := regexp.MustCompile(tomIDFormatDNS)
	reNamespace := regexp.MustCompile(tomIDNamespDNS)
	reShort := regexp.MustCompile(tomIDShortDNS)
	s = strings.TrimSuffix(s, `.`)

	switch {
	case reCommon.MatchString(s):
		sn := reCommon.FindStringSubmatch(s)
		return sn[1], sn[2], sn[3]
	case reShort.MatchString(s):
		sn := reShort.FindStringSubmatch(s)
		return sn[1], sn[2], nttShort2Long(sn[3])
	case reNamespace.MatchString(s):
		sn := reNamespace.FindStringSubmatch(s)
		return sn[1], ``, sn[2]
	default:
		return ``, ``, ``
	}

}

func parseTomIDFormatURI(s string) (name, namespace, entity string) {
	reCommon := regexp.MustCompile(tomIDFormatURI)
	reNamespace := regexp.MustCompile(tomIDNamespURI)

	switch {
	case reCommon.MatchString(s):
		sn := reCommon.FindStringSubmatch(s)
		return sn[3], sn[1], sn[2]
	case reNamespace.MatchString(s):
		sn := reNamespace.FindStringSubmatch(s)
		return sn[2], ``, sn[1]
	default:
		return ``, ``, ``
	}
}

func ParseTomIDWildcard(s string) (name, namespace, entity string) {
	reQueryNTT := regexp.MustCompile(tomIDQueryNTT)
	reSQueryNTT := regexp.MustCompile(tomIDSQueryNTT)
	reQueryNsNTT := regexp.MustCompile(tomIDQueryNsNTT)
	reSQueryNsNTT := regexp.MustCompile(tomIDSQueryNsNTT)

	switch {
	case reQueryNTT.MatchString(s):
		sn := reQueryNTT.FindStringSubmatch(s)
		return ``, ``, sn[1]
	case reSQueryNTT.MatchString(s):
		sn := reSQueryNTT.FindStringSubmatch(s)
		return ``, ``, nttShort2Long(sn[1])
	case reQueryNsNTT.MatchString(s):
		sn := reQueryNsNTT.FindStringSubmatch(s)
		return ``, sn[1], sn[2]
	case reSQueryNsNTT.MatchString(s):
		sn := reSQueryNsNTT.FindStringSubmatch(s)
		return ``, sn[1], nttShort2Long(sn[2])
	default:
		return ``, ``, ``
	}
}

func assessTomID(entity, value string) error {
	if entity != value {
		return ErrInvalidTomID
	}
	return nil
}

func nttShort2Long(s string) string {
	switch s {
	case nttContainerShort:
		return EntityContainer
	case nttOrchestrationShort:
		return EntityOrchestration
	case nttRuntimeShort:
		return EntityRuntime
	case nttServerShort:
		return EntityServer
	case nttSocketShort:
		return EntitySocket
	default:
		return ``
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
