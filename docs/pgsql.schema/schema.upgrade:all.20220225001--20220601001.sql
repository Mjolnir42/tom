BEGIN;
  CREATE SCHEMA IF NOT EXISTS abstract;
  CREATE SCHEMA IF NOT EXISTS production;
  ALTER TABLE ix.deployment_group           SET SCHEMA production;
  ALTER TABLE ix.endpoint                   SET SCHEMA production;
  
  ALTER TABLE ix.technical_service                                   SET SCHEMA production;
  ALTER TABLE ix.technical_service_standard_attribute_values         SET SCHEMA production;
  ALTER TABLE ix.technical_service_unique_attribute_values           SET SCHEMA production;
  ALTER TABLE production.technical_service                           RENAME TO instance;
  ALTER TABLE production.technical_service_standard_attribute_values RENAME TO instance_standard_attribute_values;
  ALTER TABLE production.technical_service_unique_attribute_values   RENAME TO instance_unique_attribute_values;

  ALTER TABLE ix.product                    RENAME TO consumer_product;

  ALTER TABLE ix.functional_component                                 SET SCHEMA abstract;
  ALTER TABLE ix.functional_component_standard_attribute_values       SET SCHEMA abstract;
  ALTER TABLE ix.functional_component_unique_attribute_values         SET SCHEMA abstract;
  ALTER TABLE abstract.functional_component                           RENAME TO module;
  ALTER TABLE abstract.functional_component_standard_attribute_values RENAME TO module_standard_attribute_values;
  ALTER TABLE abstract.functional_component_unique_attribute_values   RENAME TO module_unique_attribute_values;


  ALTER TABLE bulk.technical_instance       RENAME TO execution;


  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'meta', 20191016001, 'use better names for attribute value assignment tables');
COMMIT;
