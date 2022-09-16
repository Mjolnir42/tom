BEGIN;
  CREATE SCHEMA IF NOT EXISTS abstract;
  CREATE SCHEMA IF NOT EXISTS production;
  SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory, abstract, production;
  ALTER DATABASE tom SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory, abstract, production;

  ALTER TABLE bulk.technical_instance RENAME TO execution;

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

  DROP  TABLE ix.deployment_group_mapping;
  ALTER TABLE ix.deployment_group                                     SET SCHEMA production;
  ALTER TABLE ix.deployment_group_standard_attribute_values           SET SCHEMA production;
  ALTER TABLE ix.deployment_group_unique_attribute_values             SET SCHEMA production;

  ALTER TABLE ix.technical_service                                    SET SCHEMA production;
  ALTER TABLE ix.technical_service_standard_attribute_values          SET SCHEMA production;
  ALTER TABLE ix.technical_service_unique_attribute_values            SET SCHEMA production;
  ALTER TABLE production.technical_service                            RENAME TO instance;
  ALTER TABLE production.technical_service_standard_attribute_values  RENAME TO instance_standard_attribute_values;
  ALTER TABLE production.technical_service_unique_attribute_values    RENAME TO instance_unique_attribute_values;

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

  DROP  TABLE ix.endpoint_mapping;
  ALTER TABLE ix.endpoint                                             SET SCHEMA production;
  ALTER TABLE ix.endpoint_standard_attribute_values                   SET SCHEMA production;
  ALTER TABLE ix.endpoint_unique_attribute_values                     SET SCHEMA production;
  -- XXX new netrange

  -- SCHEMA abstract
  -- XXX new blueprint
  DROP  VIEW  view.deployment_group_details;
  DROP  TABLE ix.functional_component_parent;
  ALTER TABLE ix.functional_component                                 SET SCHEMA abstract;
  ALTER TABLE ix.functional_component_standard_attribute_values       SET SCHEMA abstract;
  ALTER TABLE ix.functional_component_unique_attribute_values         SET SCHEMA abstract;
  ALTER TABLE abstract.functional_component                           RENAME TO module;
  ALTER TABLE abstract.functional_component_standard_attribute_values RENAME TO module_standard_attribute_values;
  ALTER TABLE abstract.functional_component_unique_attribute_values   RENAME TO module_unique_attribute_values;
  -- XXX new abstract.artifact
  -- XXX new abstract.data
  -- XXX new abstract.service

  -- SCHEMA iX
  DROP  TABLE ix.product_mapping;
  ALTER TABLE ix.product                                              RENAME TO consumer_product;
  ALTER TABLE ix.product_standard_attribute_values                    RENAME TO consumer_product_standard_attribute_values;
  ALTER TABLE ix.product_unique_attribute_values                      RENAME TO consumer_product_unique_attribute_values;

  DROP  TABLE ix.top_level_service_mapping;








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
