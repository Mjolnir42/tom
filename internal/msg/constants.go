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
	CategoryMeta     = proto.ModelMeta
	SectionNamespace = proto.EntityNamespace
)

const (
	CategoryAsset        = proto.ModelAsset
	SectionContainer     = proto.EntityContainer
	SectionOrchestration = proto.EntityOrchestration
	SectionRuntime       = proto.EntityRuntime
	SectionServer        = proto.EntityServer
	SectionSocket        = proto.EntitySocket
)

const (
	CategoryIAM    = proto.ModelIAM
	SectionLibrary = proto.EntityLibrary
	SectionTeam    = proto.EntityTeam
	SectionUser    = proto.EntityUser
)

const (
	ActionAdd        = proto.ActionAdd
	ActionAttrAdd    = proto.ActionAttrAdd
	ActionAttrRemove = proto.ActionAttrRemove
	ActionHdSet      = proto.ActionHdSet
	ActionHdUnset    = proto.ActionHdUnset
	ActionList       = proto.ActionList
	ActionMbrAdd     = proto.ActionMbrAdd
	ActionMbrList    = proto.ActionMbrList
	ActionMbrRemove  = proto.ActionMbrRemove
	ActionMbrSet     = proto.ActionMbrSet
	ActionPropSet    = proto.ActionPropSet
	ActionPropUpdate = proto.ActionPropUpdate
	ActionRemove     = proto.ActionRemove
	ActionShow       = proto.ActionShow
	ActionUpdate     = proto.ActionUpdate
)

const (
	// RFC3339Milli is a format string for millisecond precision RFC3339
	RFC3339Milli string = `2006-01-02T15:04:05.000Z07:00`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
