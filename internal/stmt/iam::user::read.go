/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	UserReadStatements = ``

	UserList = `
SELECT      'TODO::UserList'::text;`  // TODO

	UserShow = `
SELECT      'TODO::UserShow'::text;`  // TODO

	UserKey = `
SELECT      publickey,
            fingerprint
FROM        inventory.user_key
JOIN        inventory.user
  ON        inventory.user_key.userID = inventory.user.userID
JOIN        inventory.identity_library
  ON        inventory.user.identityLibraryID = inventory.identity_library.identityLibraryID
WHERE       $3::timestamptz(3) <@ inventory.user_key.validity
  AND       inventory.user.uID = $1::text
  AND       inventory.identity_library.name = $2::text
  AND       inventory.user.isActive
  AND NOT   inventory.user.isDeleted;`
)

func init() {
	m[UserList] = `UserList`
	m[UserShow] = `UserShow`
	m[UserKey] = `UserKey`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
