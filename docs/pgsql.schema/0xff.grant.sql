--
--
-- GRANT Permissions
GRANT   SELECT
   ON   ALL TABLES IN SCHEMA view, public
   TO   tomsvc;
GRANT   SELECT,
        INSERT,
        UPDATE,
        DELETE
   ON   ALL TABLES IN SCHEMA asset, bulk, filter, ix, meta, yp
   TO   tomsvc;
GRANT   USAGE,
        SELECT
   ON   ALL SEQUENCES IN SCHEMA asset, bulk, filter, ix, meta, yp
   TO   tomsvc;
--
--
-- CREATE INITIAL BOOTSTRAP USER
BEGIN;
  SET CONSTRAINTS ALL DEFERRED;

  INSERT INTO inventory.identity_library (
    identityLibraryID,
    name,
    createdBy
  ) VALUES (
    '00000000-0000-0000-0000-000000000000'::uuid,
    'system',
    '00000000-0000-0000-0000-000000000000'::uuid,
  );
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
    '00000000-0000-0000-0000-000000000000'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    'Charlie',
    'Root',
    'root',
    0,
    'devnull@example.invalid',
    '00000000-0000-0000-0000-000000000000'::uuid
  );
  INSERT INTO inventory.team (
    teamID,
    identityLibraryID,
    name,
    createdBy
  ) VALUES (
    '00000000-0000-0000-0000-000000000000'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    'wheel',
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
    '00000000-0000-0000-0000-000000000000'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    '[-infinity,infinity]'::tstzrange,
    '00000000-0000-0000-0000-000000000000'::uuid
  );
  INSERT INTO inventory.team_lead (
    identityLibraryID,
    userID,
    headOf,
    validity,
    createdBy
  ) VALUES (
    '00000000-0000-0000-0000-000000000000'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    '00000000-0000-0000-0000-000000000000'::uuid,
    '[-infinity,infinity]'::tstzrange,
    '00000000-0000-0000-0000-000000000000'::uuid
  );
COMMIT;
