--
--
-- iX SCHEMA
CREATE TABLE IF NOT EXISTS ix.functional_component_parent (
    componentID                   uuid        NOT NULL,
    groupID                       uuid        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_ixmfc_compID  FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixmfc_grpID   FOREIGN KEY ( groupID ) REFERENCES ix.deployment_group ( groupID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __ixmfc_temporal   EXCLUDE     USING gist (public.uuid_to_bytea(groupID) WITH =,
                                                          validity WITH &&)
);
