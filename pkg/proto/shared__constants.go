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
	ActionList       = `list`
	ActionMbrAdd     = `member.add`
	ActionMbrList    = `member.list`
	ActionMbrRemove  = `member.remove`
	ActionMbrSet     = `member.set`
	ActionPropRemove = `property.remove`
	ActionPropSet    = `property.set`
	ActionPropUpdate = `property.update`
	ActionRemove     = `remove`
	ActionShow       = `show`
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
	CmdNamespace           = ModelMeta + `::` + EntityNamespace + `:`
	CmdNamespaceAdd        = ModelMeta + `::` + EntityNamespace + `:` + ActionAdd
	CmdNamespaceAttrAdd    = ModelMeta + `::` + EntityNamespace + `:` + ActionAttrAdd
	CmdNamespaceAttrRemove = ModelMeta + `::` + EntityNamespace + `:` + ActionAttrRemove
	CmdNamespaceList       = ModelMeta + `::` + EntityNamespace + `:` + ActionList
	CmdNamespacePropRemove = ModelMeta + `::` + EntityNamespace + `:` + ActionPropRemove
	CmdNamespacePropSet    = ModelMeta + `::` + EntityNamespace + `:` + ActionPropSet
	CmdNamespacePropUpdate = ModelMeta + `::` + EntityNamespace + `:` + ActionPropUpdate
	CmdNamespaceRemove     = ModelMeta + `::` + EntityNamespace + `:` + ActionRemove
	CmdNamespaceShow       = ModelMeta + `::` + EntityNamespace + `:` + ActionShow

	CmdServer           = ModelAsset + `::` + EntityServer + `:`
	CmdServerAdd        = ModelAsset + `::` + EntityServer + `:` + ActionAdd
	CmdServerList       = ModelAsset + `::` + EntityServer + `:` + ActionList
	CmdServerPropRemove = ModelAsset + `::` + EntityServer + `:` + ActionPropRemove
	CmdServerPropSet    = ModelAsset + `::` + EntityServer + `:` + ActionPropSet
	CmdServerPropUpdate = ModelAsset + `::` + EntityServer + `:` + ActionPropUpdate
	CmdServerRemove     = ModelAsset + `::` + EntityServer + `:` + ActionRemove
	CmdServerShow       = ModelAsset + `::` + EntityServer + `:` + ActionShow

	CmdRuntime           = ModelAsset + `::` + EntityRuntime + `:`
	CmdRuntimeAdd        = ModelAsset + `::` + EntityRuntime + `:` + ActionAdd
	CmdRuntimeList       = ModelAsset + `::` + EntityRuntime + `:` + ActionList
	CmdRuntimePropRemove = ModelAsset + `::` + EntityRuntime + `:` + ActionPropRemove
	CmdRuntimePropSet    = ModelAsset + `::` + EntityRuntime + `:` + ActionPropSet
	CmdRuntimePropUpdate = ModelAsset + `::` + EntityRuntime + `:` + ActionPropUpdate
	CmdRuntimeRemove     = ModelAsset + `::` + EntityRuntime + `:` + ActionRemove
	CmdRuntimeShow       = ModelAsset + `::` + EntityRuntime + `:` + ActionShow
)

const (
	TemplateCommand = `command`
	TemplateList    = `list`
	TemplateDetail  = `detail`
)

const (
	PlHoldNone  = ``
	PlHoldTomID = `:tomID`
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
