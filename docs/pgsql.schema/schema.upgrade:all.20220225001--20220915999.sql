BEGIN;
  CREATE SCHEMA IF NOT EXISTS abstract;
  CREATE SCHEMA IF NOT EXISTS production;
  SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory, abstract, production;
  ALTER DATABASE tom SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory, abstract, production;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'abstract',   20220915999, 'modelupdate' ),
                     ( 'production', 20220915999, 'modelupdate' );
COMMIT;
