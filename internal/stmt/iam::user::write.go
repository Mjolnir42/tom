/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/tom/internal/stmt"

const (
	UserWriteStatements = ``

	UserAdd = `
SELECT      'TODO::UserAdd'::text;`  // TODO

	UserRemove = `
SELECT      'TODO::UserRemove'::text;`  // TODO

	UserUpdate = `
SELECT      'TODO::UserUpdate'::text;`  // TODO

	MachineEnrol = `
INSERT INTO       inventory.user (
                         identityLibraryID,
                         firstName,
                         lastName,
                         uid,
                         isActive,
                         createdBy
                  )
SELECT            $1::uuid,
                  $2::text,
                  $3::text,
                  $4::text,
                  'yes'::boolean,
                  $5::uuid
RETURNING         userID;`

	MachineUpdateUID = `
UPDATE            inventory.user
   SET            createdBy = $1::uuid
WHERE             inventory.user.userID            = $1::uuid
  AND             inventory.user.identityLibraryID = $2::uuid;`

	UserAddKey = `
WITH sel_usr AS ( SELECT inventory.user.userID
                  FROM   inventory.user
                  JOIN   inventory.identity_library
                    ON   inventory.identity_library.identityLibraryID
                     =   inventory.user.identityLibraryID
                  WHERE  inventory.user.uid              = $4::text
                    AND  inventory.identity_library.name = $5::text )
INSERT INTO       inventory.user_key (
                         userID,
                         publicKey,
                         fingerprint,
                         createdBy
                  )
SELECT            $1::uuid,
                  $2::text,
                  $3::text,
                  $6::uuid;`

	UserActivate = `
UPDATE            inventory.user
   SET            isActive = 'yes'::boolean
FROM              inventory.identity_library
WHERE             inventory.user.identityLibraryID = inventory.identity_library.identityLibraryID
  AND             inventory.user.uID = $1::text
  AND             inventory.identity_library.name = $2::text
  AND NOT         inventory.user.isDeleted;`
)

func init() {
	m[MachineEnrol] = `MachineEnrol`
	m[MachineUpdateUID] = `MachineUpdateUID`
	m[UserActivate] = `UserActivate`
	m[UserAddKey] = `UserAddKey`
	m[UserAdd] = `UserAdd`
	m[UserRemove] = `UserRemove`
	m[UserUpdate] = `UserUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
