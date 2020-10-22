BEGIN;
  ALTER TABLE meta.dictionary_attr_values RENAME TO dictionary_standard_attribute_values;
  ALTER TABLE meta.dictionary_attr_uniq_values RENAME TO dictionary_unique_attribute_values;

  ALTER TABLE ix.top_level_service_attr_values RENAME TO top_level_service_standard_attribute_values;
  ALTER TABLE ix.top_level_service_attr_uniq_values RENAME TO top_level_service_unique_attribute_values;
  ALTER TABLE ix.product_attr_values RENAME TO product_standard_attribute_values;
  ALTER TABLE ix.product_attr_uniq_values RENAME TO product_unique_attribute_values;
  ALTER TABLE ix.information_system_attr_values RENAME TO information_system_standard_attribute_values;
  ALTER TABLE ix.information_system_attr_uniq_values RENAME TO information_system_unique_attribute_values;
  ALTER TABLE ix.functional_component_attr_values RENAME TO functional_component_standard_attribute_values;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME TO functional_component_unique_attribute_values;
  ALTER TABLE ix.deployment_group_attr_values RENAME TO deployment_group_standard_attribute_values;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME TO deployment_group_unique_attribute_values;
  ALTER TABLE ix.technical_system_service_attr_values RENAME TO technical_system_service_standard_attribute_values;
  ALTER TABLE ix.technical_system_service_attr_uniq_values RENAME TO technical_system_service_unique_attribute_values;

  ALTER TABLE asset.orchestration_environment_attribute RENAME TO orchestration_environment_standard_attribute_values;
  ALTER TABLE asset.orchestration_environment_unique RENAME TO orchestration_environment_unique_attribute_values;
  ALTER TABLE asset.runtime_environment_attribute RENAME TO runtime_environment_standard_attribute_values;
  ALTER TABLE asset.runtime_environment_unique RENAME TO runtime_environment_unique_attribute_values;
  ALTER TABLE asset.server_attribute RENAME TO server_standard_attribute_values;
  ALTER TABLE asset.server_unique RENAME TO server_unique_attribute_values;

  ALTER TABLE filter.authenticity_attribute RENAME TO authenticity_standard_attribute_values;
  ALTER TABLE filter.authenticity_unique RENAME TO authenticity_unique_attribute_values;
  ALTER TABLE filter.availability_attribute RENAME TO availability_standard_attribute_values;
  ALTER TABLE filter.availability_unique RENAME TO availability_unique_attribute_values;
  ALTER TABLE filter.brand_attribute RENAME TO brand_standard_attribute_values;
  ALTER TABLE filter.brand_unique RENAME TO brand_unique_attribute_values;
  ALTER TABLE filter.confidentiality_attribute RENAME TO confidentiality_standard_attribute_values;
  ALTER TABLE filter.confidentiality_unique RENAME TO confidentiality_unique_attribute_values;
  ALTER TABLE filter.family_attribute RENAME TO family_standard_attribute_values;
  ALTER TABLE filter.family_unique RENAME TO family_unique_attribute_values;
  ALTER TABLE filter.integrity_attribute RENAME TO integrity_standard_attribute_values;
  ALTER TABLE filter.integrity_unique RENAME TO integrity_unique_attribute_values;
  ALTER TABLE filter.lifecycle_attribute RENAME TO lifecycle_standard_attribute_values;
  ALTER TABLE filter.lifecycle_unique RENAME TO lifecycle_unique_attribute_values;
  ALTER TABLE filter.product_unit_attribute RENAME TO product_unit_standard_attribute_values;
  ALTER TABLE filter.product_unit_unique RENAME TO product_unit_unique_attribute_values;
  ALTER TABLE filter.responsible_attribute RENAME TO responsible_standard_attribute_values;
  ALTER TABLE filter.responsible_unique RENAME TO responsible_unique_attribute_values;
  ALTER TABLE filter.runner_attribute RENAME TO runner_standard_attribute_values;
  ALTER TABLE filter.runner_unique RENAME TO runner_unique_attribute_values;
  ALTER TABLE filter.service_tower_attribute RENAME TO service_tower_standard_attribute_values;
  ALTER TABLE filter.service_tower_unique RENAME TO service_tower_unique_attribute_values;
  ALTER TABLE filter.tenant_attribute RENAME TO tenant_standard_attribute_values;
  ALTER TABLE filter.tenant_unique RENAME TO tenant_unique_attribute_values;

  ALTER TABLE tosm.corporate_domain_attribute RENAME TO corporate_domain_standard_attribute_values;
  ALTER TABLE tosm.corporate_domain_unique RENAME TO corporate_domain_unique_attribute_values;
  ALTER TABLE tosm.domain_attribute RENAME TO domain_standard_attribute_values;
  ALTER TABLE tosm.domain_unique RENAME TO domain_unique_attribute_values;

  ALTER TABLE bulk.service_instances RENAME TO service_deployment_instance;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'ix', 20191017001, 'adpopt long value table names'),
                     ( 'meta', 20191017001, 'adpopt long value table names'),
                     ( 'asset', 20191017001, 'adpopt long value table names'),
                     ( 'filter', 20191017001, 'adpopt long value table names'),
                     ( 'tosm', 20191017001, 'adpopt long value table names'),
                     ( 'bulk', 20191017001, 'align table naming');
COMMIT;
