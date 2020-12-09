/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"github.com/mjolnir42/tom/pkg/proto"
)

const (
	CategoryMeta     = `meta`
	SectionNamespace = proto.EntityNamespace
)

const (
	CategoryAsset        = `asset`
	SectionContainer     = proto.EntityContainer
	SectionOrchestration = proto.EntityOrchestration
	SectionRuntime       = proto.EntityRuntime
	SectionServer        = proto.EntityServer
	SectionSocket        = proto.EntitySocket
)

const (
	ActionAdd    = `add`
	ActionList   = `list`
	ActionRemove = `remove`
	ActionShow   = `show`
	ActionAttrAdd    = `attribute.add`
	ActionAttrRemove = `attribute.remove`
	ActionPropSet    = `property.set`
	ActionPropUpdate = `property.update`
)

const (
	// RFC3339Milli is a format string for millisecond precision RFC3339
	RFC3339Milli string = `2006-01-02T15:04:05.000Z07:00`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
