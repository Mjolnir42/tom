BEGIN;
  CREATE SCHEMA IF NOT EXISTS abstract;
  CREATE SCHEMA IF NOT EXISTS production;
  SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory, abstract, production;
  ALTER DATABASE tom SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory, abstract, production;

  DROP  FUNCTION view.deployment_group_details_at;
  DROP  FUNCTION view.functional_component_details_at;
  DROP  FUNCTION view.information_system_details_at;
  DROP  VIEW     view.deployment_group_details;
  DROP  VIEW     view.functional_component_details;
  DROP  VIEW     view.information_system_details;
  DROP  TABLE    ix.deployment_group_mapping;
  DROP  TABLE    ix.endpoint_mapping;
  DROP  TABLE    ix.functional_component_parent;
  DROP  TABLE    ix.product_mapping;
  DROP  TABLE    ix.top_level_service_mapping;
  DROP  TABLE    yp.information_system_parent;
  DROP  TABLE    yp.service_mapping;

  -- SCHEMA: abstract
  -- abstract BLUEPRINT
  CREATE TABLE IF NOT EXISTS abstract.blueprint (
      bpoID                         uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_abpo          PRIMARY KEY     ( bpoID ),
      CONSTRAINT __fk_abpo_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __abpo_fk_origin   UNIQUE          ( bpoID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS abstract.blueprint_standard_attribute_values (
      bpoID                         uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_abpoav_bpoID  FOREIGN KEY     ( bpoID ) REFERENCES abstract.blueprint ( bpoID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_abpoav_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_abpoav_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_abpoav_uq_dct FOREIGN KEY     ( bpoID, dictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
      CONSTRAINT __fk_abpoav_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __abpoav_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(bpoID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS abstract.blueprint_unique_attribute_values (
      bpoID                         uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_abpoqv_bpoID  FOREIGN KEY     ( bpoID ) REFERENCES abstract.blueprint ( bpoID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_abpoqv_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_abpoqv_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_abpoqv_uq_dct FOREIGN KEY     ( bpoID, dictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
      CONSTRAINT __fk_abpoqv_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __abpoqv_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(bpoID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __abpoqv_temp_uniq EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );
  -- abstract MODULE
  ALTER TABLE ix.functional_component                                 SET SCHEMA abstract;
  ALTER TABLE ix.functional_component_standard_attribute_values       SET SCHEMA abstract;
  ALTER TABLE ix.functional_component_unique_attribute_values         SET SCHEMA abstract;
  ALTER TABLE abstract.functional_component                           RENAME TO module;
  ALTER TABLE abstract.functional_component_standard_attribute_values RENAME TO module_standard_attribute_values;
  ALTER TABLE abstract.functional_component_unique_attribute_values   RENAME TO module_unique_attribute_values;
  ALTER TABLE abstract.module                                         RENAME COLUMN componentID TO moduleID;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME COLUMN componentID TO moduleID;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME COLUMN componentID TO moduleID;
  ALTER TABLE abstract.module                                         RENAME CONSTRAINT __pk_ixfc TO __pk_amod;
  ALTER TABLE abstract.module                                         RENAME CONSTRAINT __fk_ixfc_dictID   TO __fk_amod_dictID;
  ALTER TABLE abstract.module                                         RENAME CONSTRAINT __ixfc_fk_origin   TO __amod_fk_origin;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME CONSTRAINT __fk_ixfcav_compID TO __fk_amodav_modID;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME CONSTRAINT __fk_ixfcav_attrID TO __fk_amodav_attrID;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME CONSTRAINT __fk_ixfcav_dictID TO __fk_amodav_dictID;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME CONSTRAINT __fk_ixfcav_uq_dct TO __fk_amodav_uq_dct;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME CONSTRAINT __fk_ixfcav_uq_att TO __fk_amodav_uq_att;
  ALTER TABLE abstract.module_standard_attribute_values               RENAME CONSTRAINT __ixfcav_temporal  TO __amodav_temporal;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __fk_ixfcqv_compID TO __fk_amodqv_modID;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __fk_ixfcqv_attrID TO __fk_amodqv_attrID;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __fk_ixfcqv_dictID TO __fk_amodqv_dictID;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __fk_ixfcqv_uq_dct TO __fk_amodqv_uq_dct;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __fk_ixfcqv_uq_att TO __fk_amodqv_uq_att;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __ixfcqv_temporal  TO __amodqv_temporal;
  ALTER TABLE abstract.module_unique_attribute_values                 RENAME CONSTRAINT __ixfcqv_temp_uniq TO __amodqv_temp_uniq;
  -- abstract ARTIFACT
  CREATE TABLE IF NOT EXISTS abstract.artifact (
      artifactID                    uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_aart          PRIMARY KEY     ( artifactID ),
      CONSTRAINT __fk_aart_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __aart_fk_origin   UNIQUE          ( artifactID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS abstract.artifact_standard_attribute_values (
      artifactID                    uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_aartav_artfID FOREIGN KEY     ( artifactID ) REFERENCES abstract.artifact ( artifactID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_aartav_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_aartav_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_aartav_uq_dct FOREIGN KEY     ( artifactID, dictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
      CONSTRAINT __fk_aartav_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __aartav_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(artifactID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS abstract.artifact_unique_attribute_values (
      artifactID                    uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_aartqv_artfID FOREIGN KEY     ( artifactID ) REFERENCES abstract.artifact ( artifactID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_aartqv_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_aartqv_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_aartqv_uq_dct FOREIGN KEY     ( artifactID, dictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
      CONSTRAINT __fk_aartqv_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __aartqv_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(artifactID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __aartqv_temp_uniq EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );
  -- abstract DATA
  CREATE TABLE IF NOT EXISTS abstract.data (
      dataID                        uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_adat          PRIMARY KEY     ( dataID ),
      CONSTRAINT __fk_adat_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __adat_fk_origin   UNIQUE          ( dataID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS abstract.data_standard_attribute_values (
      dataID                        uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_adatav_dataID FOREIGN KEY     ( dataID ) REFERENCES abstract.data ( dataID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_adatav_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_adatav_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_adatav_uq_dct FOREIGN KEY     ( dataID, dictionaryID ) REFERENCES abstract.data ( dataID, dictionaryID ),
      CONSTRAINT __fk_adatav_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __adatav_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(dataID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS abstract.data_unique_attribute_values (
      dataID                        uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_adatqv_dataID FOREIGN KEY     ( dataID ) REFERENCES abstract.data ( dataID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_adatqv_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_adatqv_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_adatqv_uq_dct FOREIGN KEY     ( dataID, dictionaryID ) REFERENCES abstract.data ( dataID, dictionaryID ),
      CONSTRAINT __fk_adatqv_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __adatqv_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(dataID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __adatqv_temp_uniq EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );
  -- abstract SERVICE
  CREATE TABLE IF NOT EXISTS abstract.service (
      serviceID                     uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_asrv          PRIMARY KEY     ( serviceID ),
      CONSTRAINT __fk_asrv_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __asrv_fk_origin   UNIQUE          ( serviceID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS abstract.service_standard_attribute_values (
      serviceID                     uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_asrvav_servID FOREIGN KEY     ( serviceID ) REFERENCES abstract.service ( serviceID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asrvav_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asrvav_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asrvav_uq_dct FOREIGN KEY     ( serviceID, dictionaryID ) REFERENCES abstract.service ( serviceID, dictionaryID ),
      CONSTRAINT __fk_asrvav_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __asrvav_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS abstract.service_unique_attribute_values (
      serviceID                     uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_asrvqv_servID FOREIGN KEY     ( serviceID ) REFERENCES abstract.service ( serviceID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asrvqv_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asrvqv_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asrvqv_uq_dct FOREIGN KEY     ( serviceID, dictionaryID ) REFERENCES abstract.service ( serviceID, dictionaryID ),
      CONSTRAINT __fk_asrvqv_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __asrvqv_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __asrvqv_temp_uniq EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );

  -- SCHEMA: production
  -- production TECHNICAL PRODUCT
  CREATE TABLE IF NOT EXISTS production.technical_product (
      tpID                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_ptp           PRIMARY KEY     ( tpID ),
      CONSTRAINT __fk_ptp_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __ptp_fk_origin    UNIQUE          ( tpID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS production.technical_product_standard_attribute_values (
      tpID                          uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_ptpa_tpID     FOREIGN KEY     ( tpID ) REFERENCES production.technical_product ( tpID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ptpa_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ptpa_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ptpa_uq_dct   FOREIGN KEY     ( tpID, dictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
      CONSTRAINT __fk_ptpa_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __ptpa_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(tpID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS production.technical_product_unique_attribute_values (
      tpID                          uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_ptpq_tpID     FOREIGN KEY     ( tpID ) REFERENCES production.technical_product ( tpID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ptpq_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ptpq_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_ptpq_uq_dct   FOREIGN KEY     ( tpID, dictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
      CONSTRAINT __fk_ptpq_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __ptpq_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(tpID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __ptpq_temp_uniq   EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );
  -- production DEPLOYMENT
  ALTER TABLE ix.deployment_group                                     SET SCHEMA production;
  ALTER TABLE ix.deployment_group_standard_attribute_values           SET SCHEMA production;
  ALTER TABLE ix.deployment_group_unique_attribute_values             SET SCHEMA production;
  ALTER TABLE production.deployment_group                             RENAME TO deployment;
  ALTER TABLE production.deployment_group_standard_attribute_values   RENAME TO deployment_standard_attribute_values;
  ALTER TABLE production.deployment_group_unique_attribute_values     RENAME TO deployment_unique_attribute_values;
  ALTER TABLE production.deployment                                   RENAME COLUMN groupID TO deployID;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME COLUMN groupID TO deployID;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME COLUMN groupID TO deployID;
  ALTER TABLE production.deployment                                   RENAME CONSTRAINT __pk_ixdg          TO __pk_pdpl;
  ALTER TABLE production.deployment                                   RENAME CONSTRAINT __fk_ixdg_dictID   TO __fk_pdpl_dictID;
  ALTER TABLE production.deployment                                   RENAME CONSTRAINT __ixdg_fk_origin   TO __pdpl_fk_origin;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME CONSTRAINT __fk_ixdgav_grpID  TO __fk_pdplav_dplID;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME CONSTRAINT __fk_ixdgav_attrID TO __fk_pdplav_attrID;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME CONSTRAINT __fk_ixdgav_dictID TO __fk_pdplav_dictID;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME CONSTRAINT __fk_ixdgav_uq_dct TO __fk_pdplav_uq_dct;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME CONSTRAINT __fk_ixdgav_uq_att TO __fk_pdplav_uq_att;
  ALTER TABLE production.deployment_standard_attribute_values         RENAME CONSTRAINT __ixdgav_temporal  TO __pdplav_temporal;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __fk_ixdgqv_grpID  TO __fk_pdplqv_dplID;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __fk_ixdgqv_attrID TO __fk_pdplqv_attrID;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __fk_ixdgqv_dictID TO __fk_pdplqv_dictID;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __fk_ixdgqv_uq_dct TO __fk_pdplqv_uq_dct;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __fk_ixdgqv_uq_att TO __fk_pdplqv_uq_att;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __ixdgqv_temporal  TO __pdplqv_temporal;
  ALTER TABLE production.deployment_unique_attribute_values           RENAME CONSTRAINT __ixdgqv_temp_uniq TO __pdplqv_temp_uniq;
  -- production INSTANCE
  ALTER TABLE ix.technical_service                                    SET SCHEMA production;
  ALTER TABLE ix.technical_service_standard_attribute_values          SET SCHEMA production;
  ALTER TABLE ix.technical_service_unique_attribute_values            SET SCHEMA production;
  ALTER TABLE production.technical_service                            RENAME TO instance;
  ALTER TABLE production.technical_service_standard_attribute_values  RENAME TO instance_standard_attribute_values;
  ALTER TABLE production.technical_service_unique_attribute_values    RENAME TO instance_unique_attribute_values;
  ALTER TABLE production.instance                                     RENAME COLUMN techsrvID TO instanceID;
  ALTER TABLE production.instance_standard_attribute_values           RENAME COLUMN techsrvID TO instanceID;
  ALTER TABLE production.instance_unique_attribute_values             RENAME COLUMN techsrvID TO instanceID;
  ALTER TABLE production.instance                                     RENAME CONSTRAINT __pk_ixtss         TO __pk_pinst;
  ALTER TABLE production.instance                                     RENAME CONSTRAINT __fk_ixtss_dictID  TO __fk_pinst_dictID;
  ALTER TABLE production.instance                                     RENAME CONSTRAINT __ixtss_fk_origin  TO __pinst_fk_origin;
  ALTER TABLE production.instance_standard_attribute_values           RENAME CONSTRAINT __fk_ixtssa_techID TO __fk_pinsta_instID;
  ALTER TABLE production.instance_standard_attribute_values           RENAME CONSTRAINT __fk_ixtssa_attrID TO __fk_pinsta_attrID;
  ALTER TABLE production.instance_standard_attribute_values           RENAME CONSTRAINT __fk_ixtssa_dictID TO __fk_pinsta_dictID;
  ALTER TABLE production.instance_standard_attribute_values           RENAME CONSTRAINT __fk_ixtssa_uq_dct TO __fk_pinsta_uq_dct;
  ALTER TABLE production.instance_standard_attribute_values           RENAME CONSTRAINT __fk_ixtssa_uq_att TO __fk_pinsta_uq_att;
  ALTER TABLE production.instance_standard_attribute_values           RENAME CONSTRAINT __ixtssa_temporal  TO __pinsta_temporal;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __fk_ixtssq_techID TO __fk_pinstq_instID;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __fk_ixtssq_attrID TO __fk_pinstq_attrID;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __fk_ixtssq_dictID TO __fk_pinstq_dictID;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __fk_ixtssq_uq_dct TO __fk_pinstq_uq_dct;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __fk_ixtssq_uq_att TO __fk_pinstq_uq_att;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __ixtssq_temporal  TO __pinstq_temporal;
  ALTER TABLE production.instance_unique_attribute_values             RENAME CONSTRAINT __ixtssq_temp_uniq TO __pinstq_temp_uniq;
  -- production SHARD
  CREATE TABLE IF NOT EXISTS production.shard (
      shID                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_psh           PRIMARY KEY     ( shID ),
      CONSTRAINT __fk_psh_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __psh_fk_origin    UNIQUE          ( shID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS production.shard_standard_attribute_values (
      shID                          uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_psha_shID     FOREIGN KEY     ( shID ) REFERENCES production.shard ( shID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_psha_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_psha_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_psha_uq_dct   FOREIGN KEY     ( shID, dictionaryID ) REFERENCES production.shard ( shID, dictionaryID ),
      CONSTRAINT __fk_psha_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __psha_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(shID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS production.shard_unique_attribute_values (
      shID                          uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_pshq_shID     FOREIGN KEY     ( shID ) REFERENCES production.shard ( shID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pshq_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pshq_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pshq_uq_dct   FOREIGN KEY     ( shID, dictionaryID ) REFERENCES production.shard ( shID, dictionaryID ),
      CONSTRAINT __fk_pshq_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __pshq_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(shID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __pshq_temp_uniq   EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );
  -- production ENDPOINT
  ALTER TABLE ix.endpoint                                             SET SCHEMA production;
  ALTER TABLE ix.endpoint_standard_attribute_values                   SET SCHEMA production;
  ALTER TABLE ix.endpoint_unique_attribute_values                     SET SCHEMA production;
  ALTER TABLE production.endpoint                                     RENAME CONSTRAINT __pk_ixep          TO __pk_pept;
  ALTER TABLE production.endpoint                                     RENAME CONSTRAINT __fk_ixep_dictID   TO __fk_pept_dictID;
  ALTER TABLE production.endpoint                                     RENAME CONSTRAINT __ixep_fk_origin   TO __pept_fk_origin;
  ALTER TABLE production.endpoint_standard_attribute_values           RENAME CONSTRAINT __fk_ixepsa_epID   TO __fk_peptsa_epID;
  ALTER TABLE production.endpoint_standard_attribute_values           RENAME CONSTRAINT __fk_ixepsa_attrID TO __fk_peptsa_attrID;
  ALTER TABLE production.endpoint_standard_attribute_values           RENAME CONSTRAINT __fk_ixepsa_dictID TO __fk_peptsa_dictID;
  ALTER TABLE production.endpoint_standard_attribute_values           RENAME CONSTRAINT __fk_ixepsa_uq_dct TO __fk_peptsa_uq_dct;
  ALTER TABLE production.endpoint_standard_attribute_values           RENAME CONSTRAINT __fk_ixepsa_uq_att TO __fk_peptsa_uq_att;
  ALTER TABLE production.endpoint_standard_attribute_values           RENAME CONSTRAINT __ixepsa_temporal  TO __peptsa_temporal;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __fk_ixepqv_epID   TO __fk_peptqv_epID;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __fk_ixepqv_attrID TO __fk_peptqv_attrID;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __fk_ixepqv_dictID TO __fk_peptqv_dictID;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __fk_ixepqv_uq_dct TO __fk_peptqv_uq_dct;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __fk_ixepqv_uq_att TO __fk_peptqv_uq_att;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __ixepqv_temporal  TO __peptqv_temporal;
  ALTER TABLE production.endpoint_unique_attribute_values             RENAME CONSTRAINT __ixepqv_temp_uniq TO __peptqv_temp_uniq;
  -- production NETRANGE
  CREATE TABLE IF NOT EXISTS production.netrange (
      rangeID                       uuid            NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __pk_pnr           PRIMARY KEY     ( rangeID ),
      CONSTRAINT __fk_pnr_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __pnr_fk_origin    UNIQUE          ( rangeID, dictionaryID )
  );
  CREATE TABLE IF NOT EXISTS production.netrange_standard_attribute_values (
      rangeID                       uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_pnra_rangeID  FOREIGN KEY     ( rangeID ) REFERENCES production.netrange ( rangeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pnra_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pnra_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pnra_uq_dct   FOREIGN KEY     ( rangeID, dictionaryID ) REFERENCES production.netrange ( rangeID, dictionaryID ),
      CONSTRAINT __fk_pnra_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __pnra_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(rangeID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&)
  );
  CREATE TABLE IF NOT EXISTS production.netrange_unique_attribute_values (
      rangeID                       uuid            NOT NULL,
      attributeID                   uuid            NOT NULL,
      dictionaryID                  uuid            NOT NULL,
      value                         text            NOT NULL,
      validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
      createdBy                     uuid            NOT NULL,
      createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
      CONSTRAINT __fk_pnrq_rangeID  FOREIGN KEY     ( rangeID ) REFERENCES production.netrange ( rangeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pnrq_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pnrq_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_pnrq_uq_dct   FOREIGN KEY     ( rangeID, dictionaryID ) REFERENCES production.netrange ( rangeID, dictionaryID ),
      CONSTRAINT __fk_pnrq_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
      CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
      CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
      CONSTRAINT __pnrq_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(rangeID) WITH =,
                                                                public.uuid_to_bytea(attributeID) WITH =,
                                                                validity WITH &&),
      CONSTRAINT __pnrq_temp_uniq   EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                                public.uuid_to_bytea(dictionaryID) WITH =,
                                                                value WITH =,
                                                                validity WITH &&)
  );


  -- SCHEMA iX
  ALTER TABLE ix.product                                              RENAME TO consumer_product;
  ALTER TABLE ix.product_standard_attribute_values                    RENAME TO consumer_product_standard_attribute_values;
  ALTER TABLE ix.product_unique_attribute_values                      RENAME TO consumer_product_unique_attribute_values;
  -- SCHEMA bulk
  ALTER TABLE    bulk.technical_instance                                 RENAME TO execution;
  ALTER TABLE    bulk.execution                                          RENAME COLUMN techsrvID TO instanceID;
  ALTER TABLE    bulk.execution                                          RENAME CONSTRAINT __fk_bktssi_techID TO __fk_bkexec_instID;
  ALTER TABLE    bulk.execution                                          RENAME CONSTRAINT __fk_bktssi_rteID  TO __fk_bkexec_rteID;
  ALTER TABLE    bulk.execution                                          RENAME CONSTRAINT __fk_bktssi_contID TO __fk_bkexec_contID;
  ALTER TABLE    bulk.execution                                          ADD COLUMN orchID uuid NULL;
  ALTER TABLE    bulk.execution                                          ADD CONSTRAINT __fk_bkexec_orchID FOREIGN KEY     ( orchID )
                                                                             REFERENCES asset.orchestration_environment ( orchID ) ON DELETE RESTRICT;
  ALTER TABLE    bulk.execution                                          DROP CONSTRAINT __bktssi_nonnull;
  ALTER TABLE    bulk.execution                                          ADD CONSTRAINT __bkexec_nonnull CHECK
                                                            (   (( rteID IS NOT NULL ) AND ( containerID IS     NULL ) AND ( orchID IS     NULL))
                                                             OR (( rteID IS     NULL ) AND ( containerID IS NOT NULL ) AND ( orchID IS     NULL))
                                                             OR (( rteID IS     NULL ) AND ( containerID IS     NULL ) AND ( orchID IS NOT NULL)));
  ALTER TABLE    bulk.execution                                          DROP CONSTRAINT __bktssi_temporal;
  ALTER TABLE    bulk.execution                                          ADD CONSTRAINT __bkexec_temporal EXCLUDE USING gist
                                                            (public.uuid_to_bytea(instanceID) WITH =,
                                                             public.uuid_to_bytea(rteID) WITH =,
                                                             public.uuid_to_bytea(containerID) WITH =,
                                                             public.uuid_to_bytea(orchID) WITH =,
                                                             activity WITH &&);


  -- SCHEMA yp
  -- XXX service_parent -> the table needs to be information_system_parent
  --                       since the service is the child





  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'abstract',   20220915999, 'modelupdate' ),
                     ( 'meta',       20220915999, 'modelupdate' ),
                     ( 'bulk',       20220915999, 'modelupdate' ),
                     ( 'inventory',  20220915999, 'modelupdate' ),
                     ( 'yp',         20220915999, 'modelupdate' ),
                     ( 'asset',      20220915999, 'modelupdate' ),
                     ( 'filter',     20220915999, 'modelupdate' ),
                     ( 'view',       20220915999, 'modelupdate' ),
                     ( 'ix',         20220915999, 'modelupdate' ),
                     ( 'production', 20220915999, 'modelupdate' );
COMMIT;
