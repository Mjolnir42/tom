BEGIN;
  ALTER TABLE inventory.user ALTER COLUMN firstname DROP NOT NULL;
  ALTER TABLE inventory.user ALTER COLUMN lastname DROP NOT NULL;
  ALTER TABLE inventory.user ALTER COLUMN employeeNumber DROP NOT NULL;
  ALTER TABLE inventory.user ALTER COLUMN mailAddress DROP NOT NULL;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'inventory', 20210727001, 'make user attributes optional' );
COMMIT;
