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
                         isActive
                  )
SELECT            $1::uuid,
                  $2::text,
                  $3::text,
                  $4::text,
                  'yes'::boolean
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
                  WHERE  inventory.user.uid             = $6::text
                    AND  inventory.identity_library.name = $7::text )
INSERT INTO       inventory.user_key (
	                     userID,
                         publicKey,
                         fingerprint,
						 validity,
						 createdBy
                  )
SELECT            $1::uuid,
                  $2::text,
                  $3::text,
                  tstzrange( $4::timestamptz(3), $5::timestamptz(3), '[]'),
                  $6::uuid;`
)

func init() {
	m[MachineEnrol] = `MachineEnrol`
	m[MachineUpdateUID] = `MachineUpdateUID`
	m[UserAddKey] = `UserAddKey`
	m[UserAdd] = `UserAdd`
	m[UserRemove] = `UserRemove`
	m[UserUpdate] = `UserUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
