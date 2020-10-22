BEGIN;
  DROP TABLE IF EXISTS filter.authenticity_unique_attribute_values;
  DROP TABLE IF EXISTS filter.availability_unique_attribute_values;
  DROP TABLE IF EXISTS filter.brand_unique_attribute_values;
  DROP TABLE IF EXISTS filter.confidentiality_unique_attribute_values;
  DROP TABLE IF EXISTS filter.family_unique_attribute_values;
  DROP TABLE IF EXISTS filter.integrity_unique_attribute_values;
  DROP TABLE IF EXISTS filter.lifecycle_unique_attribute_values;
  DROP TABLE IF EXISTS filter.product_unit_unique_attribute_values;
  DROP TABLE IF EXISTS filter.responsible_unique_attribute_values;
  DROP TABLE IF EXISTS filter.runner_unique_attribute_values;
  DROP TABLE IF EXISTS filter.service_tower_unique_attribute_values;
  DROP TABLE IF EXISTS filter.tenant_unique_attribute_values;
  DROP TABLE IF EXISTS filter.authenticity_standard_attribute_values;
  DROP TABLE IF EXISTS filter.availability_standard_attribute_values;
  DROP TABLE IF EXISTS filter.brand_standard_attribute_values;
  DROP TABLE IF EXISTS filter.confidentiality_standard_attribute_values;
  DROP TABLE IF EXISTS filter.family_standard_attribute_values;
  DROP TABLE IF EXISTS filter.integrity_standard_attribute_values;
  DROP TABLE IF EXISTS filter.lifecycle_standard_attribute_values;
  DROP TABLE IF EXISTS filter.product_unit_standard_attribute_values;
  DROP TABLE IF EXISTS filter.responsible_standard_attribute_values;
  DROP TABLE IF EXISTS filter.runner_standard_attribute_values;
  DROP TABLE IF EXISTS filter.service_tower_standard_attribute_values;
  DROP TABLE IF EXISTS filter.tenant_standard_attribute_values;
  DROP TABLE IF EXISTS filter.authenticity_mapping;
  DROP TABLE IF EXISTS filter.availability_mapping;
  DROP TABLE IF EXISTS filter.brand_mapping;
  DROP TABLE IF EXISTS filter.confidentiality_mapping;
  DROP TABLE IF EXISTS filter.family_mapping;
  DROP TABLE IF EXISTS filter.integrity_mapping;
  DROP TABLE IF EXISTS filter.lifecycle_mapping;
  DROP TABLE IF EXISTS filter.product_unit_mapping;
  DROP TABLE IF EXISTS filter.responsible_mapping;
  DROP TABLE IF EXISTS filter.runner_mapping;
  DROP TABLE IF EXISTS filter.service_tower_mapping;
  DROP TABLE IF EXISTS filter.tenant_mapping;
  DROP TABLE IF EXISTS filter.authenticity;
  DROP TABLE IF EXISTS filter.availability;
  DROP TABLE IF EXISTS filter.brand;
  DROP TABLE IF EXISTS filter.confidentiality;
  DROP TABLE IF EXISTS filter.family;
  DROP TABLE IF EXISTS filter.integrity;
  DROP TABLE IF EXISTS filter.lifecycle;
  DROP TABLE IF EXISTS filter.product_unit;
  DROP TABLE IF EXISTS filter.responsible;
  DROP TABLE IF EXISTS filter.runner;
  DROP TABLE IF EXISTS filter.service_tower;
  DROP TABLE IF EXISTS filter.tenant;

  CREATE TYPE flt_card AS ENUM(
      'one',
      'many'
  );
  CREATE TYPE flt_aggr AS ENUM(
      'min',
      'max',
      'first',
      'last'
  );
  CREATE TYPE flt_ntt AS ENUM(
      'top_level_service',
      'product',
      'information_system',
      'functional_component',
      'deployment_group',
      'runtime_environment',
      'orchestration_environment',
      'server'
  );

  CREATE TABLE IF NOT EXISTS filter.filter (
      filterID                      uuid        NOT NULL DEFAULT public.gen_random_uuid(),
      dictionaryID                  uuid        NOT NULL,
      CONSTRAINT __pk_ff            PRIMARY KEY ( filterID ),
      CONSTRAINT __fk_ff_dictID     FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
      CONSTRAINT __ff_fk_origin     UNIQUE      ( filterID, dictionaryID )
  );

  CREATE TABLE IF NOT EXISTS filter.name (
      filterID                      uuid        NOT NULL,
      dictionaryID                  uuid        NOT NULL,
      name                          text        NOT NULL,
      cardinality                   flt_card    NOT NULL DEFAULT 'one'::flt_card,
      aggregation                   flt_aggr    NOT NULL DEFAULT 'max'::flt_aggr,
      CONSTRAINT __fk_ffn_origin    FOREIGN KEY ( filterID, dictionaryID ) REFERENCES filter.filter ( filterID, dictionaryID ) DEFERRABLE,
      CONSTRAINT __ffn_fk_card      UNIQUE      ( filterID, cardinality ),
      CONSTRAINT __ffn_uniq_name    UNIQUE      ( dictionaryID, name )
  );

  CREATE TABLE IF NOT EXISTS filter.value (
      filterValueID                 uuid        NOT NULL DEFAULT public.gen_random_uuid(),
      filterID                      uuid        NOT NULL,
      value                         text        NOT NULL,
      valueOrder                    smallint    NOT NULL DEFAULT 0,
      CONSTRAINT __pk_ffv           PRIMARY KEY ( filterValueID ),
      CONSTRAINT __fk_ffv_filterID  FOREIGN KEY ( filterID ) REFERENCES filter.filter ( filterID ) DEFERRABLE,
      CONSTRAINT __ffv_uniq_value   UNIQUE      ( filterID, value ),
      CONSTRAINT __ffv_fk_origin    UNIQUE      ( filterValueID, filterID )
  );

  CREATE TABLE IF NOT EXISTS filter.assignable_entity (
      filterID                      uuid        NOT NULL,
      entity                        flt_ntt     NOT NULL,
      CONSTRAINT __fk_ffae_filterID FOREIGN KEY ( filterID ) REFERENCES filter.filter ( filterID ) DEFERRABLE,
      CONSTRAINT __ffae_fk_origin   UNIQUE      ( filterID, entity )
  );

CREATE TABLE IF NOT EXISTS filter.value_assignment__one (
    filterValueID                 uuid        NOT NULL,
    filterID                      uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    cardinality                   flt_card    NOT NULL,
    entity                        flt_ntt     NOT NULL,
    tlsID                         uuid        NULL,
    productID                     uuid        NULL,
    isID                          uuid        NULL,
    componentID                   uuid        NULL,
    groupID                       uuid        NULL,
    orchID                        uuid        NULL,
    rteID                         uuid        NULL,
    serverID                      uuid        NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_fvao_value    FOREIGN KEY ( filterValueID, filterID ) REFERENCES filter.value ( filterValueID, filterID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_filter   FOREIGN KEY ( filterID, dictionaryID ) REFERENCES filter.filter ( filterID, dictionaryID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_card     FOREIGN KEY ( filterID, cardinality ) REFERENCES filter.name ( filterID, cardinality ) DEFERRABLE,
    CONSTRAINT __fk_fvao_assign   FOREIGN KEY ( filterID, entity ) REFERENCES filter.assignable_entity ( filterID, entity ) DEFERRABLE,
    CONSTRAINT __fk_fvao_tlsID    FOREIGN KEY ( tlsID ) REFERENCES ix.top_level_service ( tlsID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_prodID   FOREIGN KEY ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_isID     FOREIGN KEY ( isID ) REFERENCES ix.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_compID   FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_groupID  FOREIGN KEY ( groupID ) REFERENCES ix.deployment_group ( groupID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_orchID   FOREIGN KEY ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_rteID    FOREIGN KEY ( rteID ) REFERENCES asset.runtime_environment ( rteID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_serverID FOREIGN KEY ( serverID ) REFERENCES asset.server ( serverID ) DEFERRABLE,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __fvao_cardinality CHECK       ( cardinality = 'one'::flt_card ),
    CONSTRAINT __fvao_uniq_object CHECK       (   ((entity='top_level_service')         AND (tlsID       IS NOT NULL))
                                               OR ((entity='product')                   AND (productID   IS NOT NULL))
                                               OR ((entity='information_system')        AND (isID        IS NOT NULL))
                                               OR ((entity='functional_component')      AND (componentID IS NOT NULL))
                                               OR ((entity='deployment_group')          AND (groupID     IS NOT NULL))
                                               OR ((entity='orchestration_environment') AND (orchID      IS NOT NULL))
                                               OR ((entity='runtime_environment')       AND (rteID       IS NOT NULL))
                                               OR ((entity='server')                    AND (serverID    IS NOT NULL))),
    CONSTRAINT __fvao_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(filterID) WITH =,
                                                          public.uuid_to_bytea(tlsID) WITH =,
                                                          public.uuid_to_bytea(productID) WITH =,
                                                          public.uuid_to_bytea(isID) WITH =,
                                                          public.uuid_to_bytea(componentID) WITH =,
                                                          public.uuid_to_bytea(groupID) WITH =,
                                                          public.uuid_to_bytea(orchID) WITH =,
                                                          public.uuid_to_bytea(rteID) WITH =,
                                                          public.uuid_to_bytea(serverID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS filter.value_assignment__many (
    filterValueID                 uuid        NOT NULL,
    filterID                      uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    cardinality                   flt_card    NOT NULL,
    entity                        flt_ntt     NOT NULL,
    tlsID                         uuid        NULL,
    productID                     uuid        NULL,
    isID                          uuid        NULL,
    componentID                   uuid        NULL,
    groupID                       uuid        NULL,
    orchID                        uuid        NULL,
    rteID                         uuid        NULL,
    serverID                      uuid        NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_fvam_value    FOREIGN KEY ( filterValueID, filterID ) REFERENCES filter.value ( filterValueID, filterID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_filter   FOREIGN KEY ( filterID, dictionaryID ) REFERENCES filter.filter ( filterID, dictionaryID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_card     FOREIGN KEY ( filterID, cardinality ) REFERENCES filter.name ( filterID, cardinality ) DEFERRABLE,
    CONSTRAINT __fk_fvam_assign   FOREIGN KEY ( filterID, entity ) REFERENCES filter.assignable_entity ( filterID, entity ) DEFERRABLE,
    CONSTRAINT __fk_fvam_tlsID    FOREIGN KEY ( tlsID ) REFERENCES ix.top_level_service ( tlsID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_prodID   FOREIGN KEY ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_isID     FOREIGN KEY ( isID ) REFERENCES ix.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_compID   FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_groupID  FOREIGN KEY ( groupID ) REFERENCES ix.deployment_group ( groupID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_orchID   FOREIGN KEY ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_rteID    FOREIGN KEY ( rteID ) REFERENCES asset.runtime_environment ( rteID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_serverID FOREIGN KEY ( serverID ) REFERENCES asset.server ( serverID ) DEFERRABLE,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __fvam_cardinality CHECK       ( cardinality = 'many'::flt_card ),
    CONSTRAINT __fvam_uniq_object CHECK       (   ((entity='top_level_service')         AND (tlsID       IS NOT NULL))
                                               OR ((entity='product')                   AND (productID   IS NOT NULL))
                                               OR ((entity='information_system')        AND (isID        IS NOT NULL))
                                               OR ((entity='functional_component')      AND (componentID IS NOT NULL))
                                               OR ((entity='deployment_group')          AND (groupID     IS NOT NULL))
                                               OR ((entity='orchestration_environment') AND (orchID      IS NOT NULL))
                                               OR ((entity='runtime_environment')       AND (rteID       IS NOT NULL))
                                               OR ((entity='server')                    AND (serverID    IS NOT NULL))),
    CONSTRAINT __fvam_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(filterValueID) WITH =,
                                                          public.uuid_to_bytea(tlsID) WITH =,
                                                          public.uuid_to_bytea(productID) WITH =,
                                                          public.uuid_to_bytea(isID) WITH =,
                                                          public.uuid_to_bytea(componentID) WITH =,
                                                          public.uuid_to_bytea(groupID) WITH =,
                                                          public.uuid_to_bytea(orchID) WITH =,
                                                          public.uuid_to_bytea(rteID) WITH =,
                                                          public.uuid_to_bytea(serverID) WITH =,
                                                          validity WITH &&)
);

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'filter', 20191126001, 'fully re-design filter schema');
COMMIT;
