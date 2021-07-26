--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.orchestration_environment_mapping (
    orchID                        uuid            NOT NULL,
    parentRuntimeID               uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __fk_aoep_orchID   FOREIGN KEY     ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoep_rtenv    FOREIGN KEY     ( parentRuntimeID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
);
CREATE TABLE IF NOT EXISTS asset.socket_mapping (
    socketID                      uuid            NOT NULL,
    socketDictionaryID            uuid            NOT NULL,
    endpointID                    uuid            NOT NULL,
    endpointDictionaryID          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __fk_asm_sockID    FOREIGN KEY     ( socketID, socketDictionaryID ) REFERENCES asset.socket ( socketID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asm_endpID    FOREIGN KEY     ( endpointID, endpointDictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __asm_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(socketID) WITH =,
                                                              public.uuid_to_bytea(endpointID) WITH =,
                                                              validity WITH &&)
);
