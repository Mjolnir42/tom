--
--
-- BULK SCHEMA
CREATE TABLE IF NOT EXISTS bulk.technical_instance (
    techsrvID                     uuid            NOT NULL,
    rteID                         uuid            NULL,
    containerID                   uuid            NULL,
    activity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_bktssi_techID FOREIGN KEY     ( techsrvID ) REFERENCES ix.technical_service ( techsrvID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_bktssi_rteID  FOREIGN KEY     ( rteID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_bktssi_contID FOREIGN KEY     ( containerID ) REFERENCES asset.container ( containerID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __activeFrom_utc   CHECK           ( EXTRACT( TIMEZONE FROM lower( activity ) ) = '0' ),
    CONSTRAINT __activeUntil_utc  CHECK           ( EXTRACT( TIMEZONE FROM upper( activity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __bktssi_nonnull   CHECK           (   (( rteID IS NOT NULL ) AND ( containerID IS     NULL ))
                                                   OR (( rteID IS     NULL ) AND ( containerID IS NOT NULL ))),
    CONSTRAINT __bktssi_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(techsrvID) WITH =,
                                                              public.uuid_to_bytea(rteID) WITH =,
                                                              activity WITH &&)
);
