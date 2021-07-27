BEGIN;
  ALTER TABLE inventory.user DROP CONSTRAINT __uniq_iu_uid;
  ALTER TABLE inventory.user ADD CONSTRAINT __uniq_iu_uid UNIQUE ( identityLibraryID, uid ) DEFERRABLE;
  ALTER TABLE inventory.team DROP CONSTRAINT __uniq_it_name;
  ALTER TABLE inventory.team ADD CONSTRAINT __uniq_it_name UNIQUE ( identityLibraryID, name ) DEFERRABLE;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'inventory', 20210727002, 'make inventory names unique per library' );
COMMIT;
