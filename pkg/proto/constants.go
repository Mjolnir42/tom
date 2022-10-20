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
	EntityMachine = `machine`
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
	ActionEnrolment  = `enrolment`
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
	CharDot         = `.`
	CharTilde       = `~`
	CharAsterisk    = `*`
	CharColon       = `:`
	CharUnreserved  = CharAlpha + CharDigit + CharTilde + CharDot + CharPunctuation
	CharNamespace   = CharAlpha + CharDigit + CharTilde + CharPunctuation
	CharAttribute   = CharUnreserved + CharColon
)

const (
	AttributeStandard = `standard`
	AttributeUnique   = `unique`
)

const (
	CredentialPassword = `password`
	CredentialPubKey   = `public-key`
	CredentialToken    = `token`
)

const (
	nttContainerShort     = `cnr`
	nttLibraryShort       = `lib`
	nttMachineShort       = `mac`
	nttOrchestrationShort = `ore`
	nttRuntimeShort       = `rte`
	nttServerShort        = `srv`
	nttSocketShort        = `sok`
	nttUserShort          = `usr`

	tomIDEntities = EntityServer + `|` + EntityRuntime + `|` + EntityOrchestration + `|` + EntityContainer + `|` + EntitySocket + `|` + EntityLibrary + `|` + EntityMachine + `|` + EntityTeam + `|` + EntityUser
	tomIDShortNTT = nttServerShort + `|` + nttRuntimeShort + `|` + nttOrchestrationShort + `|` + nttContainerShort + `|` + nttSocketShort + `|` + nttLibraryShort + `|` + nttMachineShort + `|` + nttUserShort

	tomIDFormatDNS = `^(?P<id>[` + CharUnreserved + `]+)\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDEntities + `)\.tom\.?$`
	tomIDShortDNS  = `^(?P<id>[` + CharUnreserved + `]+)\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDShortNTT + `)\.tom\.?$`
	tomIDNamespDNS = `^(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + EntityNamespace + `)\.tom\.?$`
	tomIDFormatURI = `^tom://(?P<ns>[` + CharNamespace + `]+)/(?P<ntt>` + tomIDEntities + `)/name=(?P<id>[` + CharUnreserved + `]+)$`
	tomIDNamespURI = `^tom:///(?P<ntt>` + EntityNamespace + `)/name=(?P<id>[` + CharNamespace + `]+)$`

	tomIDQueryNTT    = `^\` + CharAsterisk + `\.(?P<ntt>` + tomIDEntities + `)\.tom\.?$`
	tomIDSQueryNTT   = `^\` + CharAsterisk + `\.(?P<ntt>` + tomIDShortNTT + `)\.tom\.?$`
	tomIDQueryNsNTT  = `^\` + CharAsterisk + `\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDEntities + `)\.tom\.?$`
	tomIDSQueryNsNTT = `^\` + CharAsterisk + `\.(?P<ns>[` + CharNamespace + `]+)\.(?P<ntt>` + tomIDShortNTT + `)\.tom\.?$`
)

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
