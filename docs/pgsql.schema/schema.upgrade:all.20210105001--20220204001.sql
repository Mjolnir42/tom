BEGIN;
	ALTER TABLE inventory.identity_library ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE inventory.team ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE inventory.team_lead ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE inventory.team_lead ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE inventory.team_membership ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE inventory.team_membership ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE inventory.user ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');

	ALTER TABLE meta.attribute ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE meta.dictionary ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE meta.dictionary_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE meta.dictionary_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE meta.dictionary_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE meta.dictionary_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE meta.standard_attribute ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE meta.unique_attribute ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');

	ALTER TABLE ix.deployment_group ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.deployment_group ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.deployment_group_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.deployment_group_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.deployment_group_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.deployment_group_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.deployment_group_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.deployment_group_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.endpoint ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.endpoint ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.endpoint_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.endpoint_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.endpoint_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.endpoint_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.endpoint_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.endpoint_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.functional_component ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.functional_component ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.functional_component_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.functional_component_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.functional_component_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.functional_component_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.functional_component_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.functional_component_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.product ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.product ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.product_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.product_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.product_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.product_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.product_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.product_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.technical_service ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.technical_service_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.technical_service_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.technical_service_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.technical_service_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.top_level_service ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.top_level_service ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.top_level_service_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.top_level_service_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.top_level_service_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.top_level_service_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE ix.top_level_service_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE ix.top_level_service_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');

	ALTER TABLE yp.corporate_domain ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.corporate_domain_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.corporate_domain_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.corporate_domain_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.corporate_domain_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.corporate_domain_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.corporate_domain_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.domain ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.domain_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.domain_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.domain_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.domain_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.domain_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.domain_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.information_system ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.information_system ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.information_system_linking ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.information_system_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.information_system_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.information_system_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.information_system_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.information_system_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.information_system_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.service ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.service_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.service_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.service_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.service_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.service_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.service_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE yp.service_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE yp.service_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');

	ALTER TABLE asset.container ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.container_linking ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.container_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.container_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.container_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.container_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.container_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.container_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.orchestration_environment ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.orchestration_environment_linking ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.orchestration_environment_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.orchestration_environment_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.orchestration_environment_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.orchestration_environment_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.orchestration_environment_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.orchestration_environment_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.runtime_environment ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.runtime_environment_linking ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.runtime_environment_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.runtime_environment_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.runtime_environment_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.runtime_environment_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.runtime_environment_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.runtime_environment_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.server ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.server_linking ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.server_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.server_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.server_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.server_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.server_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.server_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.socket ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.socket_linking ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.socket_mapping ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.socket_mapping ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.socket_parent ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.socket_parent ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.socket_standard_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.socket_standard_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE asset.socket_unique_attribute_values ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE asset.socket_unique_attribute_values ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');

	ALTER TABLE bulk.technical_instance ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE bulk.technical_instance ALTER COLUMN activity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');

	ALTER TABLE filter.assignable_entity ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE filter.filter ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE filter.name ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE filter.value ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE filter.value_assignment__many ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE filter.value_assignment__many ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');
	ALTER TABLE filter.value_assignment__one ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');
	ALTER TABLE filter.value_assignment__one ALTER COLUMN validity SET DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]');

	ALTER TABLE public.schema_versions RENAME COLUMN created_at TO createdAt;
	ALTER TABLE public.schema_versions ALTER COLUMN createdAt SET DEFAULT (now() at time zone 'utc');

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'asset', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'bulk', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'filter', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'inventory', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'ix', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'meta', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'public', 20220204001, 'update default value for timestamp columns with timezone'),
                     ( 'yp', 20220204001, 'update default value for timestamp columns with timezone');
COMMIT;
