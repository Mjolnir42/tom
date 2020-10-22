BEGIN;
  ALTER TABLE yp.service_linking ALTER COLUMN serviceID DROP DEFAULT;
  ALTER TABLE ix.mapping_information_system SET SCHEMA yp;

  ALTER TABLE asset.orchestration_environment_parent RENAME TO orchestration_environment_mapping;
  ALTER TABLE ix.mapping_deployment_group RENAME TO deployment_group_mapping;
  ALTER TABLE ix.mapping_endpoint RENAME TO endpoint_mapping;
  ALTER TABLE ix.mapping_functional_component RENAME TO functional_component_parent;
  ALTER TABLE ix.mapping_product RENAME TO product_mapping;
  ALTER TABLE ix.mapping_top_level_service RENAME TO top_level_service_mapping;
  ALTER TABLE yp.mapping_corporate_domain RENAME TO corporate_domain_parent;
  ALTER TABLE yp.mapping_domain RENAME TO domain_parent;
  ALTER TABLE yp.mapping_information_system RENAME TO information_system_parent;

  CREATE TABLE IF NOT EXISTS asset.socket_mapping (
      socketID                      uuid        NOT NULL,
      socketDictionaryID            uuid        NOT NULL,
      endpointID                    uuid        NOT NULL,
      endpointDictionaryID          uuid        NOT NULL,
      validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      CONSTRAINT __fk_asm_sockID    FOREIGN KEY ( socketID, socketDictionaryID ) REFERENCES asset.socket ( socketID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asm_endpID    FOREIGN KEY ( endpointID, endpointDictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __asm_temporal     EXCLUDE     using gist (public.uuid_to_bytea(socketID) WITH =,
                                                            public.uuid_to_bytea(endpointID) WITH =,
                                                            validity WITH &&)
  );

  CREATE TABLE IF NOT EXISTS yp.information_system_linking (
      isLinkID                      uuid        NOT NULL DEFAULT public.gen_random_uuid(),
      isID_A                        uuid        NOT NULL,
      dictionaryID_A                uuid        NOT NULL,
      isID_B                        uuid        NOT NULL,
      dictionaryID_B                uuid        NOT NULL,
      CONSTRAINT __pk_ypisl         PRIMARY KEY ( isLinkID ),
      CONSTRAINT __fk_ypisl_isA     FOREIGN KEY ( isID_A, dictionaryID_A ) REFERENCES yp.information_system ( isID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ypisl_isB     FOREIGN KEY ( isID_B, dictionaryID_B ) REFERENCES yp.information_system ( isID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __ypisl_diff_isID  CHECK       ( isID_A != isID_B ),
      CONSTRAINT __ypisl_uniq_link  UNIQUE      ( isID_A, isID_B ),
      CONSTRAINT __ypisl_ordered    CHECK       ( public.uuid_to_bytea(isID_A) > public.uuid_to_bytea(isID_B) )
  );

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'yp', 20200917001, 'cleanup of relationship tables' ),
                     ( 'asset', 20200917001, 'cleanup of relationship tables' ),
                     ( 'ix', 20200917001, 'cleanup of relationship tables' );
COMMIT;
