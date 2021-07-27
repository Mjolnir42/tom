BEGIN;
  ALTER TABLE bulk.technical_instance ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE bulk.technical_instance ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE bulk.technical_instance ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE bulk.technical_instance ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE filter.filter ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE filter.filter ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE filter.filter ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE filter.filter ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE filter.name ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE filter.name ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE filter.name ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE filter.name ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE filter.value ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE filter.value ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE filter.value ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE filter.value ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE filter.assignable_entity ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE filter.assignable_entity ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE filter.assignable_entity ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE filter.assignable_entity ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE filter.value_assignment__one ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE filter.value_assignment__one ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE filter.value_assignment__one ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE filter.value_assignment__one ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE filter.value_assignment__many ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE filter.value_assignment__many ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE filter.value_assignment__many ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE filter.value_assignment__many ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  SAVEPOINT tables;

  ALTER TABLE bulk.technical_instance ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE filter.filter ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE filter.name ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE filter.value ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE filter.assignable_entity ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE filter.value_assignment__one ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE filter.value_assignment__many ALTER COLUMN createdBy DROP DEFAULT;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'bulk', 20210727001, 'add inventory information' ),
                     ( 'filter', 20210727001, 'add inventory information' );
COMMIT;
