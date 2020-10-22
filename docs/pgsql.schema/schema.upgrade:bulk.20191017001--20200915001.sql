BEGIN;
  ALTER TABLE bulk.service_deployment_instance RENAME TO technical_instance;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'bulk', 20200915001, 'rename instance data table');
COMMIT;
