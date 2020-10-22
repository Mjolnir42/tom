--
-- YP SCHEMA
CREATE TABLE IF NOT EXISTS yp.service (
    serviceID                     uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid        NOT NULL,
    CONSTRAINT __pk_yps           PRIMARY KEY ( serviceID ),
    CONSTRAINT __fk_yps_dictID    FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __yps_fk_origin    UNIQUE      ( serviceID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS yp.service_standard_attribute_values (
    serviceID                     uuid        NOT NULL,
    attributeID                   uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    value                         text        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_ypsa_serID    FOREIGN KEY ( serviceID ) REFERENCES yp.service ( serviceID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsa_attrID   FOREIGN KEY ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsa_dictID   FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsa_uq_dict  FOREIGN KEY ( serviceID, dictionaryID ) REFERENCES yp.service ( serviceID, dictionaryID ),
    CONSTRAINT __fk_ypsa_uq_att   FOREIGN KEY ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __ypsa_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS yp.service_unique_attribute_values (
    serviceID                     uuid        NOT NULL,
    attributeID                   uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    value                         text        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_ypsq_serID    FOREIGN KEY ( serviceID ) REFERENCES yp.service ( serviceID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsq_attrID   FOREIGN KEY ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsq_dictID   FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsq_uq_dict  FOREIGN KEY ( serviceID, dictionaryID ) REFERENCES yp.service ( serviceID, dictionaryID ),
    CONSTRAINT __fk_ypsq_uq_att   FOREIGN KEY ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __ypsq_temporal    EXCLUDE     USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&),
    CONSTRAINT __ypsq_temp_uniq   EXCLUDE     USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                          public.uuid_to_bytea(dictionaryID) WITH =,
                                                          value WITH =,
                                                          validity WITH &&)
);
