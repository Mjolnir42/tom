--
--
-- BULK SCHEMA
CREATE TABLE IF NOT EXISTS bulk.execution (
    instanceID                    uuid            NOT NULL,
    rteID                         uuid            NULL,
    containerID                   uuid            NULL,
    orchID                        uuid            NULL,
    activity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_bkexec_instID FOREIGN KEY     ( instanceID ) REFERENCES production.instance ( instanceID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_bkexec_rteID  FOREIGN KEY     ( rteID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_bkexec_contID FOREIGN KEY     ( containerID ) REFERENCES asset.container ( containerID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_bkexec_orchID FOREIGN KEY     ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __activeFrom_utc   CHECK           ( EXTRACT( TIMEZONE FROM lower( activity ) ) = '0' ),
    CONSTRAINT __activeUntil_utc  CHECK           ( EXTRACT( TIMEZONE FROM upper( activity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __bkexec_nonnull   CHECK           (   (( rteID IS NOT NULL ) AND ( containerID IS     NULL ) AND ( orchID IS     NULL))
                                                   OR (( rteID IS     NULL ) AND ( containerID IS NOT NULL ) AND ( orchID IS     NULL))
                                                   OR (( rteID IS     NULL ) AND ( containerID IS     NULL ) AND ( orchID IS NOT NULL))),
    CONSTRAINT __bkexec_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(instanceID) WITH =,
                                                              public.uuid_to_bytea(rteID) WITH =,
                                                              public.uuid_to_bytea(containerID) WITH =,
                                                              public.uuid_to_bytea(orchID) WITH =,
                                                              activity WITH &&)
);
