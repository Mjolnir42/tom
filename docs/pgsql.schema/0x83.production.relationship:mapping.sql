---
---
---
CREATE TABLE IF NOT EXISTS production.technical_product_mapping (
    tpID                          uuid            NOT NULL,
    tpDictionaryID                uuid            NOT NULL,
    deployID                      uuid            NOT NULL,
    deployDictionaryID            uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_ptpm_tpID     FOREIGN KEY     ( tpID, tpDictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
    CONSTRAINT __fk_ptpm_dplID    FOREIGN KEY     ( deployID, deployDictionaryID ) REFERENCES production.deployment ( deployID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ptpm_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(tpID) WITH =,
                                                              public.uuid_to_bytea(deployID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS production.deployment_mapping (
    deployID                      uuid            NOT NULL,
    deployDictionaryID            uuid            NOT NULL,
    instanceID                    uuid            NOT NULL,
    instanceDictionaryID          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_pdm_dplID     FOREIGN KEY     ( deployID, deployDictionaryID ) REFERENCES production.deployment ( deployID, dictionaryID ),
    CONSTRAINT __fk_pdm_insID     FOREIGN KEY     ( instanceID, instanceDictionaryID ) REFERENCES production.instance ( instanceID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __pdm_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(deployID) WITH =,
                                                              public.uuid_to_bytea(instanceID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS production.netrange_mapping (
    rangeID                       uuid            NOT NULL,
    rangeDictionaryID             uuid            NOT NULL,
    tpID                          uuid            NULL,
    tpDictionaryID                uuid            NULL,
    deployID                      uuid            NULL,
    deployDictionaryID            uuid            NULL,
    instanceID                    uuid            NULL,
    instanceDictionaryID          uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_pnm_rangeID   FOREIGN KEY     ( rangeID, rangeDictionaryID ) REFERENCES production.netrange ( rangeID, dictionaryID ),
    CONSTRAINT __fk_pnm_tpID      FOREIGN KEY     ( tpID, tpDictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
    CONSTRAINT __fk_pnm_dplID     FOREIGN KEY     ( deployID, deployDictionaryID ) REFERENCES production.deployment ( deployID, dictionaryID ),
    CONSTRAINT __fk_pnm_insID     FOREIGN KEY     ( instanceID, instanceDictionaryID ) REFERENCES production.instance ( instanceID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __pnm_nonnull      CHECK           (   ((tpID IS NOT NULL) AND (deployID IS     NULL) AND (instanceID IS     NULL))
                                                   OR ((tpID IS     NULL) AND (deployID IS NOT NULL) AND (instanceID IS     NULL))
                                                   OR ((tpID IS     NULL) AND (deployID IS     NULL) AND (instanceID IS NOT NULL))),
    CONSTRAINT __pnm_null_tp      CHECK           (   ((tpID IS NOT NULL) AND (tpDictionaryID IS NOT NULL))
                                                   OR ((tpID IS     NULL) AND (tpDictionaryID IS     NULL))),
    CONSTRAINT __pnm_null_dpl     CHECK           (   ((deployID IS NOT NULL) AND (deployDictionaryID IS NOT NULL))
                                                   OR ((deployID IS     NULL) AND (deployDictionaryID IS     NULL))),
    CONSTRAINT __pnm_null_ins     CHECK           (   ((instanceID IS NOT NULL) AND (instanceDictionaryID IS NOT NULL))
                                                   OR ((instanceID IS     NULL) AND (instanceDictionaryID IS     NULL))),
    CONSTRAINT __pnm_temp_tp      EXCLUDE         USING gist (public.uuid_to_bytea(rangeID) WITH =,
                                                              public.uuid_to_bytea(tpID) WITH =,
                                                              validity WITH &&) WHERE (tpID IS NOT NULL),
    CONSTRAINT __pnm_temp_deploy  EXCLUDE         USING gist (public.uuid_to_bytea(rangeID) WITH =,
                                                              public.uuid_to_bytea(deployID) WITH =,
                                                              validity WITH &&) WHERE (deployID IS NOT NULL),
    CONSTRAINT __pnm_temp_ins     EXCLUDE         USING gist (public.uuid_to_bytea(rangeID) WITH =,
                                                              public.uuid_to_bytea(instanceID) WITH =,
                                                              validity WITH &&) WHERE (instanceID IS NOT NULL)
);
