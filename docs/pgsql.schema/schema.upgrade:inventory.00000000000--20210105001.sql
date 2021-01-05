BEGIN;
  SET CONSTRAINTS ALL DEFERRED;
  CREATE SCHEMA IF NOT EXISTS inventory;

  SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory;
  ALTER DATABASE tom SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory;

  CREATE TABLE IF NOT EXISTS inventory.identity_library (
      identityLibraryID             uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      name                          varchar(128)    NOT NULL,
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
      CONSTRAINT __pk_iil           PRIMARY KEY     ( identityLibraryID ),
      CONSTRAINT __uniq_iil_name    UNIQUE          ( name ) DEFERRABLE,
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
  );
  CREATE TABLE IF NOT EXISTS inventory.user (
      userID                        uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      identityLibraryID             uuid            NOT NULL,
      firstName                     varchar(256)    NOT NULL,
      lastName                      varchar(256)    NOT NULL,
      uid                           varchar(256)    NOT NULL,
      employeeNumber                numeric(16,0)   NOT NULL,
      mailAddress                   text            NOT NULL,
      externalID                    text            NULL,
      isActive                      boolean         NOT NULL DEFAULT 'no',
      isDeleted                     boolean         NOT NULL DEFAULT 'no',
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
      CONSTRAINT __pk_iu            PRIMARY KEY     ( userID ),
      CONSTRAINT __fk_iu_idLibID    FOREIGN KEY     ( identityLibraryID ) REFERENCES inventory.identity_library ( identityLibraryID ) DEFERRABLE,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ) DEFERRABLE,
      CONSTRAINT __uniq_iu_empNum   UNIQUE          ( identityLibraryID, employeeNumber ) DEFERRABLE,
      CONSTRAINT __uniq_iu_extID    UNIQUE          ( identityLibraryID, externalID ) DEFERRABLE,
      CONSTRAINT __uniq_iu_uid      UNIQUE          ( uid ) DEFERRABLE,
      CONSTRAINT __iu_fk_origin     UNIQUE          ( identityLibraryID, userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
  );
  ALTER TABLE inventory.identity_library ADD CONSTRAINT
      __fk_createdBy                FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ) DEFERRABLE
  ;
  CREATE TABLE IF NOT EXISTS inventory.team (
      teamID                        uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      identityLibraryID             uuid            NOT NULL,
      name                          varchar(384)    NOT NULL,
      externalID                    text            NULL,
      isDeleted                     boolean         NOT NULL DEFAULT 'no',
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
      CONSTRAINT __pk_it            PRIMARY KEY     ( teamID ),
      CONSTRAINT __fk_it_idLibID    FOREIGN KEY     ( identityLibraryID ) REFERENCES inventory.identity_library ( identityLibraryID ) DEFERRABLE,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ) DEFERRABLE,
      CONSTRAINT __uniq_it_name     UNIQUE          ( name ) DEFERRABLE,
      CONSTRAINT __uniq_it_extID    UNIQUE          ( identityLibraryID, externalID ) DEFERRABLE,
      CONSTRAINT __it_fk_origin     UNIQUE          ( identityLibraryID, teamID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
  );
  CREATE TABLE IF NOT EXISTS inventory.team_membership (
      identityLibraryID             uuid            NOT NULL,
      userID                        uuid            NOT NULL,
      teamID                        uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
      CONSTRAINT __fk_itm_userID    FOREIGN KEY     ( identityLibraryID, userID ) REFERENCES inventory.user ( identityLibraryID, userID ) DEFERRABLE,
      CONSTRAINT __fk_itm_teamID    FOREIGN KEY     ( identityLibraryID, teamID ) REFERENCES inventory.team ( identityLibraryID, teamID ) DEFERRABLE,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ) DEFERRABLE,
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __itm_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(userID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS inventory.team_lead (
      identityLibraryID             uuid            NOT NULL,
      userID                        uuid            NOT NULL,
      headOf                        uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
      CONSTRAINT __fk_itl_userID    FOREIGN KEY     ( identityLibraryID, userID ) REFERENCES inventory.user ( identityLibraryID, userID ) DEFERRABLE,
      CONSTRAINT __fk_itl_headOf    FOREIGN KEY     ( identityLibraryID, headOf ) REFERENCES inventory.team ( identityLibraryID, teamID ) DEFERRABLE,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ) DEFERRABLE,
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __itl_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(headOf) WITH =,
                                                                validity WITH &&)
  );

  INSERT INTO inventory.identity_library (
    identityLibraryID,
    name,
    createdBy
  ) VALUES (
    '00000000-0000-0000-0000-000000000000'::uuid,
    'system',
    '00000000-0000-0000-0000-000000000000'::uuid
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

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'inventory', 20210105001, 'add new namespace: inventory' );
COMMIT;
