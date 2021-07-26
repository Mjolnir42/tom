BEGIN;
  ALTER TABLE ix.top_level_service ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.top_level_service ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.top_level_service ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.top_level_service ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.top_level_service_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.top_level_service_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.top_level_service_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.top_level_service_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.top_level_service_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.top_level_service_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.top_level_service_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.top_level_service_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.top_level_service_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.top_level_service_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.top_level_service_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.top_level_service_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.product ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.product ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.product ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.product ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.product_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.product_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.product_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.product_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.product_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.product_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.product_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.product_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.product_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.product_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.product_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.product_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.information_system ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.information_system ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.information_system ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.information_system ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.information_system_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.information_system_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.information_system_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.information_system_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.information_system_linking ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.information_system_linking ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.information_system_linking ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.information_system_linking ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.information_system_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.information_system_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.information_system_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.information_system_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.information_system_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.information_system_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.information_system_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.information_system_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.functional_component ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.functional_component ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.functional_component ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.functional_component ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.functional_component_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.functional_component_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.functional_component_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.functional_component_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.functional_component_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.functional_component_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.functional_component_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.functional_component_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.functional_component_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.functional_component_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.functional_component_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.functional_component_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.deployment_group ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.deployment_group ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.deployment_group ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.deployment_group ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.deployment_group_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.deployment_group_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.deployment_group_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.deployment_group_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.deployment_group_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.deployment_group_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.deployment_group_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.deployment_group_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.deployment_group_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.deployment_group_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.deployment_group_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.deployment_group_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.technical_service ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.technical_service ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.technical_service ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.technical_service ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.technical_service_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.technical_service_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.technical_service_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.technical_service_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.technical_service_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.technical_service_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.technical_service_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.technical_service_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.endpoint ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.endpoint ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.endpoint ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.endpoint ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.endpoint_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.endpoint_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.endpoint_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.endpoint_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.endpoint_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.endpoint_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.endpoint_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.endpoint_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE ix.endpoint_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE ix.endpoint_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE ix.endpoint_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE ix.endpoint_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.corporate_domain ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.corporate_domain ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.corporate_domain ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.corporate_domain ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.corporate_domain_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.corporate_domain_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.corporate_domain_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.corporate_domain_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.corporate_domain_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.corporate_domain_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.corporate_domain_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.corporate_domain_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.corporate_domain_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.corporate_domain_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.corporate_domain_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.corporate_domain_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.domain ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.domain ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.domain ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.domain ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.domain_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.domain_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.domain_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.domain_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.domain_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.domain_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.domain_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.domain_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.domain_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.domain_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.domain_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.domain_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.service ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.service ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.service ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.service ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.service_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.service_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.service_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.service_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.service_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.service_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.service_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.service_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.service_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.service_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.service_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.service_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE yp.service_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE yp.service_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE yp.service_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE yp.service_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.server ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.server ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.server ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.server ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.server_linking ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.server_linking ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.server_linking ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.server_linking ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.server_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.server_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.server_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.server_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.server_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.server_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.server_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.server_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.server_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.server_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.server_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.server_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.runtime_environment ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.runtime_environment ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.runtime_environment ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.runtime_environment ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.runtime_environment_linking ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.runtime_environment_linking ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.runtime_environment_linking ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.runtime_environment_linking ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.runtime_environment_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.runtime_environment_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.runtime_environment_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.runtime_environment_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.runtime_environment_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.runtime_environment_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.runtime_environment_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.runtime_environment_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.runtime_environment_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.runtime_environment_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.runtime_environment_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.runtime_environment_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.orchestration_environment ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.orchestration_environment ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.orchestration_environment ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.orchestration_environment ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.orchestration_environment_linking ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.orchestration_environment_linking ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.orchestration_environment_linking ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.orchestration_environment_linking ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.orchestration_environment_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.orchestration_environment_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.orchestration_environment_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.orchestration_environment_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.orchestration_environment_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.orchestration_environment_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.orchestration_environment_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.orchestration_environment_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.orchestration_environment_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.orchestration_environment_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.orchestration_environment_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.orchestration_environment_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.socket ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.socket ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.socket ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.socket ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.socket_linking ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.socket_linking ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.socket_linking ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.socket_linking ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.socket_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.socket_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.socket_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.socket_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.socket_mapping ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.socket_mapping ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.socket_mapping ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.socket_mapping ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.socket_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.socket_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.socket_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.socket_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.socket_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.socket_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.socket_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.socket_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.container ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.container ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.container ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.container ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.container_linking ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.container_linking ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.container_linking ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.container_linking ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.container_parent ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.container_parent ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.container_parent ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.container_parent ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.container_standard_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.container_standard_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.container_standard_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.container_standard_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );

  ALTER TABLE asset.container_unique_attribute_values ADD COLUMN createdBy uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'::uuid;
  ALTER TABLE asset.container_unique_attribute_values ADD COLUMN createdAt timestamptz(3) NOT NULL DEFAULT now();
  ALTER TABLE asset.container_unique_attribute_values ADD CONSTRAINT __fk_createdBy FOREIGN KEY ( createdBy ) REFERENCES inventory.user ( userID );
  ALTER TABLE asset.container_unique_attribute_values ADD CONSTRAINT __createdAt_utc CHECK ( EXTRACT( TIMEZONE FROM createdAt ) = '0' );
  SAVEPOINT tables;

  ALTER TABLE ix.top_level_service ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.top_level_service_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.top_level_service_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.top_level_service_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.product ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.product_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.product_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.product_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.information_system ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.information_system_linking ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.information_system_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.information_system_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.information_system_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.functional_component ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.functional_component_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.functional_component_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.functional_component_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.deployment_group ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.deployment_group_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.deployment_group_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.deployment_group_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.technical_service ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.technical_service_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.technical_service_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.endpoint ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.endpoint_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.endpoint_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE ix.endpoint_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.corporate_domain ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.corporate_domain_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.corporate_domain_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.corporate_domain_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.domain ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.domain_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.domain_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.domain_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.service ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.service_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.service_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.service_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE yp.service_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.server ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.server_linking ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.server_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.server_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.server_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.runtime_environment ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.runtime_environment_linking ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.runtime_environment_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.runtime_environment_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.runtime_environment_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.orchestration_environment ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.orchestration_environment_linking ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.orchestration_environment_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.orchestration_environment_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.orchestration_environment_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.socket ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.socket_linking ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.socket_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.socket_mapping ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.socket_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.socket_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.container ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.container_linking ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.container_parent ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.container_standard_attribute_values ALTER COLUMN createdBy DROP DEFAULT;
  ALTER TABLE asset.container_unique_attribute_values ALTER COLUMN createdBy DROP DEFAULT;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'yp', 20210726001, 'add inventory information' ),
                     ( 'asset', 20210726001, 'add inventory information' ),
                     ( 'ix', 20210726001, 'add inventory information' );
COMMIT;
