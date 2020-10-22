BEGIN;
  CREATE TABLE IF NOT EXISTS ix.endpoint (
      endpointID                    uuid        NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid        NOT NULL,
      validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      CONSTRAINT __pk_ixep          PRIMARY KEY ( endpointID ),
      CONSTRAINT __fk_ixep_dictID   FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __ixep_fk_origin   UNIQUE      ( endpointID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS ix.endpoint_standard_attribute_values (
      endpointID                    uuid        NOT NULL,
      attributeID                   uuid        NOT NULL,
      dictionaryID                  uuid        NOT NULL,
      value                         text        NOT NULL,
      validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      CONSTRAINT __fk_ixepsa_epID   FOREIGN KEY ( endpointID ) REFERENCES ix.endpoint ( endpointID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixepsa_attrID FOREIGN KEY ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixepsa_dictID FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixepsa_uq_dct FOREIGN KEY ( endpointID, dictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ),
      CONSTRAINT __fk_ixepsa_uq_att FOREIGN KEY ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __ixepsa_temporal  EXCLUDE     USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                            public.uuid_to_bytea(attributeID) WITH =,
                                                            validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS ix.endpoint_unique_attribute_values (
      endpointID                    uuid        NOT NULL,
      attributeID                   uuid        NOT NULL,
      dictionaryID                  uuid        NOT NULL,
      value                         text        NOT NULL,
      validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      CONSTRAINT __fk_ixepqv_epID   FOREIGN KEY ( endpointID ) REFERENCES ix.endpoint ( endpointID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixepqv_attrID FOREIGN KEY ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixepqv_dictID FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixepqv_uq_dct FOREIGN KEY ( endpointID, dictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ),
      CONSTRAINT __fk_ixepqv_uq_att FOREIGN KEY ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __ixepqv_temporal  EXCLUDE     USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                            public.uuid_to_bytea(attributeID) WITH =,
                                                            validity WITH &&),
      CONSTRAINT __ixepqv_temp_uniq EXCLUDE     USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                            public.uuid_to_bytea(dictionaryID) WITH =,
                                                            value WITH =,
                                                            validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS ix.mapping_endpoint (
      endpointID                    uuid        NOT NULL,
      dictionaryID                  uuid        NOT NULL,
      componentID                   uuid        NULL,
      componentDictionaryID         uuid        NULL,
      groupID                       uuid        NULL,
      groupDictionaryID             uuid        NULL,
      validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      CONSTRAINT __fk_ixme_endpID   FOREIGN KEY ( endpointID ) REFERENCES ix.endpoint ( endpointID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_dictID   FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_uq_endp  FOREIGN KEY ( endpointID, dictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_cmpID    FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_cDictID  FOREIGN KEY ( componentDictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_uq_comp  FOREIGN KEY ( componentID, componentDictionaryID ) REFERENCES ix.functional_component ( componentID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_grpID    FOREIGN KEY ( groupID ) REFERENCES ix.deployment_group ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_gDictID  FOREIGN KEY ( groupDictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ixme_uq_group FOREIGN KEY ( groupID, groupDictionaryID ) REFERENCES ix.deployment_group ( groupID, dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __ixme_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                            public.uuid_to_bytea(componentID) WITH =,
                                                            public.uuid_to_bytea(groupID) WITH =,
                                                            validity WITH &&),
      CONSTRAINT __ixme_uniq_map    CHECK       (   ((componentID IS NOT NULL) AND (componentDictionaryID IS NOT NULL) AND (groupID IS     NULL) AND (groupDictionaryID IS     NULL))
                                                 OR ((componentID IS     NULL) AND (componentDictionaryID IS     NULL) AND (groupID IS NOT NULL) AND (groupDictionaryID IS NOT NULL)))
  );

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'ix', 20200914001, 'add endpoint tables');
COMMIT;

