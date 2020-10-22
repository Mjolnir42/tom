BEGIN;
  ALTER TABLE ix.information_system_attribute RENAME TO information_system_attr_values;
  ALTER TABLE ix.information_system_unique RENAME TO information_system_attr_uniq_values;

  ALTER TABLE ix.technical_system_service_attribute RENAME TO technical_system_service_attr_values;
  ALTER TABLE ix.technical_system_service_unique RENAME TO technical_system_service_attr_uniq_values;

  ALTER TABLE ix.top_level_service_attribute RENAME TO top_level_service_attr_values;
  ALTER TABLE ix.top_level_service_unique RENAME TO top_level_service_attr_uniq_values;

  ALTER TABLE ix.product_attribute RENAME TO product_attr_values;
  ALTER TABLE ix.product_unique RENAME TO product_attr_uniq_values;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'ix', 20191015001, 'use better names for attribute value assignment tables');
COMMIT;
