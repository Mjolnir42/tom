BEGIN;

  INSERT INTO inventory.user (
    userID,
    identityLibraryID,
    firstName,
    lastName,
    uid,
    employeeNumber,
    mailAddress,
    createdBy
  ) VALUES (
    'ffffffff-ffff-ffff-ffff-ffffffffffff'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    'Anonymous',
    'Unconfigured',
    'nobody',
    null,
    'devzero@example.invalid',
    '00000000-0000-0000-0000-000000000000'::uuid
  );
  INSERT INTO inventory.team_membership (
    identityLibraryID,
    userID,
    teamID,
    validity,
    createdBy
  ) VALUES (
    '00000000-0000-0000-0000-000000000000'::uuid,
    'ffffffff-ffff-ffff-ffff-ffffffffffff'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    '[-infinity,infinity]'::tstzrange,
    '00000000-0000-0000-0000-000000000000'::uuid
  );

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'inventory', 20210824001, 'add user system~nobody' );
COMMIT;
