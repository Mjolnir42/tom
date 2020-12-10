--
-- YP SCHEMA
CREATE TABLE IF NOT EXISTS yp.service_mapping (
    serviceID                     uuid        NOT NULL,
    serviceDictionaryID           uuid        NOT NULL,
    endpointID                    uuid        NOT NULL,
    endpointDictionaryID          uuid        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __fk_ypsm_serID    FOREIGN KEY ( serviceID, serviceDictionaryID ) REFERENCES yp.service ( serviceID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsm_enpID    FOREIGN KEY ( endpointID, endpointDictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __ypsm_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                          public.uuid_to_bytea(endpointID) WITH =,
                                                          validity WITH &&)
);
