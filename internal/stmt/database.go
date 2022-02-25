/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	DatabaseStatements = ``

	DatabaseTimezone = `SET TIME ZONE 'UTC';`

	DatabaseIsolationLevel = `SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;`

	ReadOnlyTransaction = `SET TRANSACTION READ ONLY, DEFERRABLE;`

	DeferredTransaction = `SET CONSTRAINTS ALL DEFERRED;`
)

func init() {
	m[DatabaseIsolationLevel] = `DatabaseIsolationLevel`
	m[DatabaseTimezone] = `DatabaseTimezone`
	m[ReadOnlyTransaction] = `ReadOnlyTransaction`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
