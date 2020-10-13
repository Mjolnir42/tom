/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

import (
	"regexp"
	"strings"
)

type Property [2]string

type AttributeDefinition struct {
	Key    string `json:"key"`
	Unique bool   `json:"uniqueValueConstraint,omitempty"`
}

const (
	tomIDFormatDNS = `^(:?[^[:space:]])\.(:?[^[:space:]]+)\.(:?server|runtime|orchestration)\.tom\.?$`
	tomIDFormatURI = `^tom://(:?[^[:space:]]+)/(:?server|runtime|orchestration)/name=(:?[^[:space:]]+)$`
)

func isTomIDFormatDNS(s string) bool {
	re := regexp.MustCompile(tomIDFormatDNS)
	return re.MatchString(s)
}

func isTomIDFormatURI(s string) bool {
	re := regexp.MustCompile(tomIDFormatURI)
	return re.MatchString(s)
}

func parseTomIDFormatDNS(s string) (name, namespace, entity string) {
	s = strings.TrimSuffix(s, `.`)
	parts := strings.Split(s, `.`)
	return strings.Join(parts[0:len(parts)-3], ``), parts[len(parts)-3], parts[len(parts)-2]
}

func parseTomIDFormatURI(s string) (name, namespace, entity string) {
	parts := strings.Split(s, `/`)
	return strings.TrimPrefix(parts[len(parts)-1], `name=`), parts[len(parts)-3], parts[len(parts)-2]
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
