/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	ModelAsset          = `asset`
	EntityContainer     = `container`
	EntityOrchestration = `orchestration`
	EntityRuntime       = `runtime`
	EntityServer        = `server`
	EntitySocket        = `socket`
)

const (
	ModelIAM      = `iam`
	EntityLibrary = `library`
	EntityTeam    = `team`
	EntityUser    = `user`
)

const (
	ModelMeta       = `meta`
	EntityNamespace = `namespace`
)

const (
	ActionAdd        = `add`
	ActionAttrAdd    = `attribute.add`
	ActionAttrRemove = `attribute.remove`
	ActionHdSet      = `headof.set`
	ActionHdUnset    = `headof.unset`
	ActionLink       = `link`
	ActionList       = `list`
	ActionMbrAdd     = `member.add`
	ActionMbrList    = `member.list`
	ActionMbrRemove  = `member.remove`
	ActionMbrSet     = `member.set`
	ActionPropRemove = `property.remove`
	ActionPropSet    = `property.set`
	ActionPropUpdate = `property.update`
	ActionRemove     = `remove`
	ActionResolve    = `resolve`
	ActionShow       = `show`
	ActionStack      = `stack`
	ActionUnstack    = `unstack`
	ActionUpdate     = `update`
)

const (
	MethodDELETE = `DELETE`
	MethodGET    = `GET`
	MethodHEAD   = `HEAD`
	MethodPATCH  = `PATCH`
	MethodPOST   = `POST`
	MethodPUT    = `PUT`
)

const (
	MetaPropertyCmdLink    = ModelAsset + `::` + `meta-cmd` + `::` + ActionLink
	MetaPropertyCmdStack   = ModelAsset + `::` + `meta-cmd` + `::` + ActionStack
	MetaPropertyCmdUnstack = ModelAsset + `::` + `meta-cmd` + `::` + ActionUnstack
)

const (
	TemplateCommand = `command`
	TemplateList    = `list`
	TemplateDetail  = `detail`
)

const (
	PlHoldNone     = ``
	PlHoldTomID    = `:tomID`
	PlHoldTargetID = `:targetID`
	PlHoldResolv   = `:level`
)

const (
	CharAlpha       = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	CharDigit       = `0123456789`
	CharPunctuation = `_-`
	CharUnreserved  = CharAlpha + CharDigit + CharPunctuation
	CharTilde       = `~`
	CharNamespace   = CharAlpha + CharDigit + CharTilde + CharPunctuation
)

const (
	AttributeStandard = `standard`
	AttributeUnique   = `unique`
)

const (
	nttContainerShort     = `cnr`
	nttOrchestrationShort = `ore`
	nttRuntimeShort       = `rte`
	nttServerShort        = `srv`
	nttSocketShort        = `sok`

	tomIDEntities = EntityServer + `|` + EntityRuntime + `|` + EntityOrchestration + `|` + EntityContainer + `|` + EntitySocket
	tomIDShortNTT = nttServerShort + `|` + nttRuntimeShort + `|` + nttOrchestrationShort + `|` + nttContainerShort + `|` + nttSocketShort

	tomIDFormatDNS = `^(?P<id>[` + CharUnreserved + `]+)\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDEntities + `)\.tom\.?$`
	tomIDShortDNS  = `^(?P<id>[` + CharUnreserved + `]+)\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDShortNTT + `)\.tom\.?$`
	tomIDNamespDNS = `^(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + EntityNamespace + `)\.tom\.?$`
	tomIDFormatURI = `^tom://(?P<ns>[` + CharNamespace + `]+)/(?P<ntt>` + tomIDEntities + `)/name=(?P<id>[` + CharUnreserved + `]+)$`
	tomIDNamespURI = `^tom:///(?P<ntt>` + EntityNamespace + `)/name=(?P<id>[` + CharNamespace + `]+)$`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
