/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// PropertyDetail holds all the information about an object property
type PropertyDetail struct {
	Attribute  string          `json:"attribute"`
	Value      string          `json:"value"`
	Raw        json.RawMessage `json:"structuredValue,omitempty"`
	ValidSince string          `json:"validSince"`
	ValidUntil string          `json:"validUntil"`
	CreatedAt  string          `json:"createdAt"`
	CreatedBy  string          `json:"createdBy"`
	Namespace  string          `json:"namespace"`
}

// AttributeDefinition holds the definition of a dictionary attribute
type AttributeDefinition struct {
	Key    string `json:"key"`
	Unique bool   `json:"uniqueValueConstraint"`
}

const (
	tomIDFormatDNS = `^(:?[` + CharUnreserved + `]+)\.(:?[` + CharNamespace + `]+)\.(:?server|runtime|orchestration|container|socket)\.tom\.?$`
	tomIDNamespDNS = `^(:?[` + CharNamespace + `]+)\.(:?namespace)\.tom\.?$`
	tomIDFormatURI = `^tom://(:?[` + CharNamespace + `]+)/(:?server|runtime|orchestration|container|socket)/name=(:?[` + CharUnreserved + `]+)$`
	tomIDNamespURI = `^tom:///(:?namespace)/name=(:?[` + CharNamespace + `]+)$`
)

func IsTomID(s string) bool {
	return isTomIDFormatDNS(s) || isTomIDFormatURI(s)
}

func isTomIDFormatDNS(s string) bool {
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", tomIDFormatDNS, tomIDNamespDNS))
	return re.MatchString(s)
}

func isTomIDFormatURI(s string) bool {
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", tomIDFormatURI, tomIDNamespURI))
	return re.MatchString(s)
}

func parseTomIDFormatDNS(s string) (name, namespace, entity string) {
	reCommon := regexp.MustCompile(tomIDFormatDNS)
	reNamespace := regexp.MustCompile(tomIDNamespDNS)
	s = strings.TrimSuffix(s, `.`)
	parts := strings.Split(s, `.`)

	switch {
	case reCommon.MatchString(s):
		return strings.Join(parts[0:len(parts)-3], ``), parts[len(parts)-3], parts[len(parts)-2]
	case reNamespace.MatchString(s):
		return strings.Join(parts[0:len(parts)-2], ``), ``, parts[len(parts)-2]
	default:
		return ``, ``, ``
	}

}

func parseTomIDFormatURI(s string) (name, namespace, entity string) {
	reCommon := regexp.MustCompile(tomIDFormatURI)
	reNamespace := regexp.MustCompile(tomIDNamespURI)

	parts := strings.Split(s, `/`)
	switch {
	case reCommon.MatchString(s):
		return strings.TrimPrefix(parts[len(parts)-1], `name=`), parts[len(parts)-3], parts[len(parts)-2]
	case reNamespace.MatchString(s):
		return strings.TrimPrefix(parts[len(parts)-1], `name=`), ``, parts[len(parts)-2]
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
