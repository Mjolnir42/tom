BEGIN;
  ALTER SCHEMA tosm RENAME TO yp;
  SET search_path TO ix, meta, filter, yp, asset, 'view', bulk;
  ALTER DATABASE ix SET search_path TO ix, meta, filter, yp, asset, 'view', bulk;

  ALTER TABLE ix.information_system SET SCHEMA yp;
  ALTER TABLE ix.information_system_standard_attribute_values SET SCHEMA yp;
  ALTER TABLE ix.information_system_unique_attribute_values SET SCHEMA yp;

  UPDATE public.schema_versions SET schema='yp' WHERE schema='tosm';
  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'yp', 20200915001, 'rename schema tosm to yp');
COMMIT;
