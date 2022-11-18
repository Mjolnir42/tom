---
---
---
CREATE TABLE IF NOT EXISTS production.shard_parent (
    shID                          uuid            NOT NULL,
    shDictionaryID                uuid            NOT NULL,
    tpID                          uuid            NULL,
    tpDictionaryID                uuid            NULL,
    deployID                      uuid            NULL,
    deployDictionaryID            uuid            NULL,
    instanceID                    uuid            NULL,
    instanceDictionaryID          uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_psm_shID      FOREIGN KEY     ( shID, shDictionaryID ) REFERENCES production.shard ( shID, dictionaryID ),
    CONSTRAINT __fk_psm_tpID      FOREIGN KEY     ( tpID, tpDictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
    CONSTRAINT __fk_psm_dplID     FOREIGN KEY     ( deployID, deployDictionaryID ) REFERENCES production.deployment ( deployID, dictionaryID ),
    CONSTRAINT __fk_psm_insID     FOREIGN KEY     ( instanceID, instanceDictionaryID ) REFERENCES production.instance ( instanceID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __psm_nonnull      CHECK           (   ((tpID IS NOT NULL) AND (deployID IS     NULL) AND (instanceID IS     NULL))
                                                   OR ((tpID IS     NULL) AND (deployID IS NOT NULL) AND (instanceID IS     NULL))
                                                   OR ((tpID IS     NULL) AND (deployID IS     NULL) AND (instanceID IS NOT NULL))),
    CONSTRAINT __psm_null_tp      CHECK           (   ((tpID IS NOT NULL) AND (tpDictionaryID IS NOT NULL))
                                                   OR ((tpID IS     NULL) AND (tpDictionaryID IS     NULL))),
    CONSTRAINT __psm_null_dpl     CHECK           (   ((deployID IS NOT NULL) AND (deployDictionaryID IS NOT NULL))
                                                   OR ((deployID IS     NULL) AND (deployDictionaryID IS     NULL))),
    CONSTRAINT __psm_null_ins     CHECK           (   ((instanceID IS NOT NULL) AND (instanceDictionaryID IS NOT NULL))
                                                   OR ((instanceID IS     NULL) AND (instanceDictionaryID IS     NULL))),
    CONSTRAINT __psm_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(shID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __psm_temp_tp      EXCLUDE         USING gist (public.uuid_to_bytea(shID) WITH =,
                                                              public.uuid_to_bytea(tpID) WITH =,
                                                              validity WITH &&) WHERE (tpID IS NOT NULL),
    CONSTRAINT __psm_temp_deploy  EXCLUDE         USING gist (public.uuid_to_bytea(shID) WITH =,
                                                              public.uuid_to_bytea(deployID) WITH =,
                                                              validity WITH &&) WHERE (deployID IS NOT NULL),
    CONSTRAINT __psm_temp_ins     EXCLUDE         USING gist (public.uuid_to_bytea(shID) WITH =,
                                                              public.uuid_to_bytea(instanceID) WITH =,
                                                              validity WITH &&) WHERE (instanceID IS NOT NULL)
);
CREATE TABLE IF NOT EXISTS production.endpoint_parent (
    endpointID                    uuid            NOT NULL,
    epDictionaryID                uuid            NOT NULL,
    tpID                          uuid            NULL,
    tpDictionaryID                uuid            NULL,
    deployID                      uuid            NULL,
    deployDictionaryID            uuid            NULL,
    instanceID                    uuid            NULL,
    instanceDictionaryID          uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_pem_endpID    FOREIGN KEY     ( endpointID, epDictionaryID ) REFERENCES production.endpoint ( endpointID, dictionaryID ),
    CONSTRAINT __fk_pem_tpID      FOREIGN KEY     ( tpID, tpDictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
    CONSTRAINT __fk_pem_dplID     FOREIGN KEY     ( deployID, deployDictionaryID ) REFERENCES production.deployment ( deployID, dictionaryID ),
    CONSTRAINT __fk_pem_insID     FOREIGN KEY     ( instanceID, instanceDictionaryID ) REFERENCES production.instance ( instanceID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __pem_nonnull      CHECK           (   ((tpID IS NOT NULL) AND (deployID IS     NULL) AND (instanceID IS     NULL))
                                                   OR ((tpID IS     NULL) AND (deployID IS NOT NULL) AND (instanceID IS     NULL))
                                                   OR ((tpID IS     NULL) AND (deployID IS     NULL) AND (instanceID IS NOT NULL))),
    CONSTRAINT __pem_null_tp      CHECK           (   ((tpID IS NOT NULL) AND (tpDictionaryID IS NOT NULL))
                                                   OR ((tpID IS     NULL) AND (tpDictionaryID IS     NULL))),
    CONSTRAINT __pem_null_dpl     CHECK           (   ((deployID IS NOT NULL) AND (deployDictionaryID IS NOT NULL))
                                                   OR ((deployID IS     NULL) AND (deployDictionaryID IS     NULL))),
    CONSTRAINT __pem_null_ins     CHECK           (   ((instanceID IS NOT NULL) AND (instanceDictionaryID IS NOT NULL))
                                                   OR ((instanceID IS     NULL) AND (instanceDictionaryID IS     NULL))),
    CONSTRAINT __pem_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __pem_temp_tp      EXCLUDE         USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                              public.uuid_to_bytea(tpID) WITH =,
                                                              validity WITH &&) WHERE (tpID IS NOT NULL),
    CONSTRAINT __pem_temp_deploy  EXCLUDE         USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                              public.uuid_to_bytea(deployID) WITH =,
                                                              validity WITH &&) WHERE (deployID IS NOT NULL),
    CONSTRAINT __pem_temp_ins     EXCLUDE         USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                              public.uuid_to_bytea(instanceID) WITH =,
                                                              validity WITH &&) WHERE (instanceID IS NOT NULL)
);
