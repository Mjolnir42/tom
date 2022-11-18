---
---
---
CREATE TABLE IF NOT EXISTS abstract.blueprint_realization (
    bpoID                         uuid            NOT NULL,
    bpoDictionaryID               uuid            NOT NULL,
    tpID                          uuid            NOT NULL,
    tpDictionaryID                uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __abr_fk_bpoID     FOREIGN KEY     ( bpoID, bpoDictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
    CONSTRAINT __abr_fk_tpID      FOREIGN KEY     ( tpID, tpDictionaryID ) REFERENCES production.technical_product ( tpID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __abr_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(tpID) WITH =,
                                                              validity WITH &&)
);
CREATE INDEX IF NOT EXISTS __abr_idx_tpID ON abstract.blueprint_realization ( tpID, bpoID );
CREATE INDEX IF NOT EXISTS __abr_idx_bpoID ON abstract.blueprint_realization ( bpoID, tpID );

CREATE TABLE IF NOT EXISTS abstract.module_realization (
    moduleID                      uuid            NOT NULL,
    moduleDictionaryID            uuid            NOT NULL,
    deployID                      uuid            NOT NULL,
    deployDictionaryID            uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __amr_fk_moduleID  FOREIGN KEY     ( moduleID, moduleDictionaryID ) REFERENCES abstract.module ( moduleID, dictionaryID ),
    CONSTRAINT __amr_fk_deployID  FOREIGN KEY     ( deployID, deployDictionaryID ) REFERENCES production.deployment ( deployID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __amr_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(deployID) WITH =,
                                                              validity WITH &&)
);
CREATE INDEX IF NOT EXISTS __amr_idx_moduleID ON abstract.module_realization ( moduleID, deployID );
CREATE INDEX IF NOT EXISTS __amr_idx_deployID ON abstract.module_realization ( deployID, moduleID );

CREATE TABLE IF NOT EXISTS abstract.artifact_realization (
    artifactID                    uuid            NOT NULL,
    artifactDictionaryID          uuid            NOT NULL,
    instanceID                    uuid            NOT NULL,
    instanceDictionaryID          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __aar_fk_artID     FOREIGN KEY     ( artifactID, artifactDictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
    CONSTRAINT __aar_fk_insID     FOREIGN KEY     ( instanceID, instanceDictionaryID ) REFERENCES production.instance ( instanceID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aar_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(instanceID) WITH =,
                                                              validity WITH &&)
);
CREATE INDEX IF NOT EXISTS __aar_idx_artifactID ON abstract.artifact_realization ( artifactID, instanceID );
CREATE INDEX IF NOT EXISTS __aar_idx_instanceID ON abstract.artifact_realization ( instanceID, artifactID );

CREATE TABLE IF NOT EXISTS abstract.data_realization (
    dataID                        uuid            NOT NULL,
    dataDictionaryID              uuid            NOT NULL,
    shardID                       uuid            NOT NULL,
    shardDictionaryID             uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __adr_fk_artID     FOREIGN KEY     ( dataID, dataDictionaryID ) REFERENCES abstract.data ( dataID, dictionaryID ),
    CONSTRAINT __adr_fk_insID     FOREIGN KEY     ( shardID, shardDictionaryID ) REFERENCES production.shard ( shID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __adr_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(shardID) WITH =,
                                                              validity WITH &&)
);
CREATE INDEX IF NOT EXISTS __adr_idx_dataID ON abstract.data_realization ( dataID, shardID );
CREATE INDEX IF NOT EXISTS __adr_idx_shardID ON abstract.data_realization ( shardID, dataID );

CREATE TABLE IF NOT EXISTS abstract.service_realization (
    serviceID                     uuid            NOT NULL,
    serviceDictionaryID           uuid            NOT NULL,
    endpointID                    uuid            NOT NULL,
    endpointDictionaryID          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __asr_fk_artID     FOREIGN KEY     ( serviceID, serviceDictionaryID ) REFERENCES abstract.service ( serviceID, dictionaryID ),
    CONSTRAINT __asr_fk_insID     FOREIGN KEY     ( endpointID, endpointDictionaryID ) REFERENCES production.endpoint ( endpointID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __asr_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                              validity WITH &&)
);
CREATE INDEX IF NOT EXISTS __asr_idx_serviceID ON abstract.service_realization ( serviceID, endpointID );
CREATE INDEX IF NOT EXISTS __asr_idx_endpointID ON abstract.service_realization ( endpointID, serviceID );

