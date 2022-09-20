IdentityLibraryID
  team
    . identitylibraryid
    . name
  team_lead
    . identitylibraryid
    . userid
    . headof (teamID)
  team_membership
    . identityLibraryID
    . userID
    . teamID
  user
    . identityLibraryID
    . userID
-----------------------------------------
  tom://${namespace}/${entity}/name=${name}
  tom://${identityLibrary}/${team|user}/name=
-----------------------------------------
  tom://engineroom/machine/uid=
-----------------------------------------
  ${user}.ims.user.tom
  ${team}.ims.team.tom
-----------------------------------------
EnrollmentKey?
  Library -> boolean self-Enrollment attribute
          -> boolean machine library attribute
          -> enrollment key attribute

  Client: create pubkey/privkey keypair

  CSR:
    - uid         => key fingerprint
    - identityLib => engineroom
    - externalID  => publicKey
    - time        => timestamp.created()

  -> fingerprint.engineroom.machine.tom

  X-TOM-INITIALIZE-ENROLLMENT= true

    runtime: pubkey
