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
	CmdOrchestration           = ModelAsset + `::` + EntityOrchestration + `:`
	CmdOrchestrationAdd        = ModelAsset + `::` + EntityOrchestration + `:` + ActionAdd
	CmdOrchestrationLink       = ModelAsset + `::` + EntityOrchestration + `:` + ActionLink
	CmdOrchestrationList       = ModelAsset + `::` + EntityOrchestration + `:` + ActionList
	CmdOrchestrationPropRemove = ModelAsset + `::` + EntityOrchestration + `:` + ActionPropRemove
	CmdOrchestrationPropSet    = ModelAsset + `::` + EntityOrchestration + `:` + ActionPropSet
	CmdOrchestrationPropUpdate = ModelAsset + `::` + EntityOrchestration + `:` + ActionPropUpdate
	CmdOrchestrationRemove     = ModelAsset + `::` + EntityOrchestration + `:` + ActionRemove
	CmdOrchestrationShow       = ModelAsset + `::` + EntityOrchestration + `:` + ActionShow
	CmdOrchestrationStack      = ModelAsset + `::` + EntityOrchestration + `:` + ActionStack
	CmdOrchestrationUnstack    = ModelAsset + `::` + EntityOrchestration + `:` + ActionUnstack

	CmdPageStaticCSS   = `page::static::css`
	CmdPageStaticFont  = `page::static::font`
	CmdPageStaticImage = `page::static::img`
	CmdPageStaticJS    = `page::static::js`
	CmdPageApplication = `page::application`

	MetaPropertyCmdLink    = ModelAsset + `::` + `meta-cmd` + `::` + ActionLink
	MetaPropertyCmdStack   = ModelAsset + `::` + `meta-cmd` + `::` + ActionStack
	MetaPropertyCmdUnstack = ModelAsset + `::` + `meta-cmd` + `::` + ActionUnstack
)

const (
	TemplateBindata = `bindata`
	TemplateCommand = `command`
	TemplateDetail  = `detail`
	TemplateList    = `list`
	TemplateNone    = `none`
)

const (
	PlHoldNone     = ``
	PlHoldTomID    = `:tomID`
	PlHoldTargetID = `:targetID`
	PlHoldAsset    = `:asset`
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
	tomIDEntities = EntityServer + `|` + EntityRuntime + `|` + EntityOrchestration + `|` + EntityContainer + `|` + EntitySocket

	tomIDFormatDNS = `^(?P<id>[` + CharUnreserved + `]+)\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDEntities + `)\.tom\.?$`
	tomIDNamespDNS = `^(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + EntityNamespace + `)\.tom\.?$`
	tomIDFormatURI = `^tom://(?P<ns>[` + CharNamespace + `]+)/(?P<ntt>` + tomIDEntities + `)/name=(?P<id>[` + CharUnreserved + `]+)$`
	tomIDNamespURI = `^tom:///(?P<ntt>` + EntityNamespace + `)/name=(?P<id>[` + CharNamespace + `]+)$`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
