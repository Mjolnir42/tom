--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.runtime_environment_parent (
    rteID                         uuid        NOT NULL,
    parentServerID                uuid        NULL,
    parentRuntimeID               uuid        NULL,
    parentOrchestrationID         uuid        NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_artep_rteID   FOREIGN KEY ( rteID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artep_server  FOREIGN KEY ( parentServerID ) REFERENCES asset.server ( serverID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artep_rtenv   FOREIGN KEY ( parentRuntimeID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artep_orch    FOREIGN KEY ( parentOrchestrationID ) REFERENCES asset.orchestration_environment ( orchID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __artep_uq_parent  CHECK       (   ((parentServerID IS NOT NULL) AND (parentRuntimeID IS     NULL) AND (parentOrchestrationID IS     NULL))
                                               OR ((parentServerID IS     NULL) AND (parentRuntimeID IS NOT NULL) AND (parentOrchestrationID IS     NULL))
                                               OR ((parentServerID IS     NULL) AND (parentRuntimeID IS     NULL) AND (parentOrchestrationID IS NOT NULL))),
    CONSTRAINT __artep_temporal   EXCLUDE     USING gist (public.uuid_to_bytea(rteID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.server_parent (
    serverID                      uuid        NOT NULL,
    parentRuntimeID               uuid        NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_asp_srvID     FOREIGN KEY ( serverID ) REFERENCES asset.server ( serverID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asp_rtenv     FOREIGN KEY ( parentRuntimeID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __asp_uq_parent    CHECK       ( parentRuntimeID IS NOT NULL ),
    CONSTRAINT __asp_temporal     EXCLUDE     USING gist (public.uuid_to_bytea(serverID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.socket_parent (
    socketID                      uuid        NOT NULL,
    parentRuntimeID               uuid        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_asop_sockID   FOREIGN KEY ( socketID ) REFERENCES asset.socket ( socketID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asop_rteID    FOREIGN KEY ( parentRuntimeID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __asop_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(socketID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.container_parent (
    containerID                   uuid        NOT NULL,
    parentRuntimeID               uuid        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_acop_sockID   FOREIGN KEY ( containerID ) REFERENCES asset.container ( containerID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_acop_rteID    FOREIGN KEY ( parentRuntimeID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __acop_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(containerID) WITH =,
                                                          validity WITH &&)
);
