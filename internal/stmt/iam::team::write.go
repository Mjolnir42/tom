/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	TeamWriteStatements = ``

	TeamAdd = `
SELECT      'TODO::TeamAdd'::text;`  // TODO

	TeamMbrAdd = `
SELECT      'TODO::TeamMbrAdd'::text;`  // TODO

	TeamMbrSet = `
SELECT      'TODO::TeamMbrSet'::text;`  // TODO

	TeamMbrRemove = `
SELECT      'TODO::TeamMbrRemove'::text;`  // TODO

	TeamHdSet = `
SELECT      'TODO::TeamHdSet'::text;`  // TODO

	TeamHdUnset = `
SELECT      'TODO::TeamHdUnset'::text;`  // TODO

	TeamRemove = `
SELECT      'TODO::TeamRemove'::text;`  // TODO

	TeamUpdate = `
SELECT      'TODO::TeamUpdate'::text;`  // TODO
)

func init() {
	m[TeamAdd] = `TeamAdd`
	m[TeamMbrAdd] = `TeamMbrAdd`
	m[TeamMbrSet] = `TeamMbrSet`
	m[TeamMbrRemove] = `TeamMbrRemove`
	m[TeamHdSet] = `TeamHdSet`
	m[TeamHdUnset] = `TeamHdUnset`
	m[TeamRemove] = `TeamRemove`
	m[TeamUpdate] = `TeamUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
