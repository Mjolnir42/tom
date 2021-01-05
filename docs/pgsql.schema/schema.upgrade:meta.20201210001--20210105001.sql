BEGIN;
  ALTER TABLE meta.dictionary ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE meta.dictionary ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE meta.dictionary ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE meta.dictionary ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE meta.attribute ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE meta.attribute ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE meta.attribute ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE meta.attribute ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE meta.standard_attribute ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE meta.standard_attribute ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE meta.standard_attribute ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE meta.standard_attribute ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE meta.unique_attribute ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE meta.unique_attribute ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE meta.unique_attribute ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE meta.unique_attribute ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE meta.dictionary_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE meta.dictionary_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE meta.dictionary_standard_attribute_values
              ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE meta.dictionary_standard_attribute_values
              ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE meta.dictionary_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE meta.dictionary_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE meta.dictionary_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE meta.dictionary_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );
  SAVEPOINT tables;

  ALTER TABLE meta.dictionary ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE meta.attribute ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE meta.standard_attribute ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE meta.unique_attribute ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE meta.dictionary_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE meta.dictionary_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'meta', 20210105001, 'add inventory information' );
COMMIT;
