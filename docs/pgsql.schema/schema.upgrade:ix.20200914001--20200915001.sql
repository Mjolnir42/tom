BEGIN;
  ALTER TABLE ix.technical_system_service RENAME TO technical_service;
  ALTER TABLE ix.technical_system_service_standard_attribute_values RENAME TO technical_service_standard_attribute_values;
  ALTER TABLE ix.technical_system_service_unique_attribute_values RENAME TO technical_service_unique_attribute_values;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'ix', 20200915001, 'rename technical system services to technical services');
COMMIT;

