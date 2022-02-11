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
	CmdServer           = ModelAsset + `::` + EntityServer + `:`
	CmdServerAdd        = ModelAsset + `::` + EntityServer + `:` + ActionAdd
	CmdServerLink       = ModelAsset + `::` + EntityServer + `:` + ActionLink
	CmdServerList       = ModelAsset + `::` + EntityServer + `:` + ActionList
	CmdServerPropRemove = ModelAsset + `::` + EntityServer + `:` + ActionPropRemove
	CmdServerPropSet    = ModelAsset + `::` + EntityServer + `:` + ActionPropSet
	CmdServerPropUpdate = ModelAsset + `::` + EntityServer + `:` + ActionPropUpdate
	CmdServerRemove     = ModelAsset + `::` + EntityServer + `:` + ActionRemove
	CmdServerShow       = ModelAsset + `::` + EntityServer + `:` + ActionShow
	CmdServerStack      = ModelAsset + `::` + EntityServer + `:` + ActionStack

	CmdRuntime           = ModelAsset + `::` + EntityRuntime + `:`
	CmdRuntimeAdd        = ModelAsset + `::` + EntityRuntime + `:` + ActionAdd
	CmdRuntimeLink       = ModelAsset + `::` + EntityRuntime + `:` + ActionLink
	CmdRuntimeList       = ModelAsset + `::` + EntityRuntime + `:` + ActionList
	CmdRuntimePropRemove = ModelAsset + `::` + EntityRuntime + `:` + ActionPropRemove
	CmdRuntimePropSet    = ModelAsset + `::` + EntityRuntime + `:` + ActionPropSet
	CmdRuntimePropUpdate = ModelAsset + `::` + EntityRuntime + `:` + ActionPropUpdate
	CmdRuntimeRemove     = ModelAsset + `::` + EntityRuntime + `:` + ActionRemove
	CmdRuntimeResolve    = ModelAsset + `::` + EntityRuntime + `:` + ActionResolve
	CmdRuntimeShow       = ModelAsset + `::` + EntityRuntime + `:` + ActionShow
	CmdRuntimeStack      = ModelAsset + `::` + EntityRuntime + `:` + ActionStack

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
	CharNamespace   = CharAlpha + CharDigit + `~` + CharPunctuation
)

const (
	AttributeStandard = `standard`
	AttributeUnique   = `unique`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
