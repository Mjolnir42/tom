--
-- YP SCHEMA
CREATE TABLE IF NOT EXISTS yp.technology_reference_card (
    trcID                         uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_yptr          PRIMARY KEY     ( trcID ),
    CONSTRAINT __fk_yptr_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __yptr_fk_origin   UNIQUE          ( trcID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS yp.technology_reference_card_standard_attribute_values (
    trcID                         uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_yptra_serID   FOREIGN KEY     ( trcID ) REFERENCES yp.technology_reference_card ( trcID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_yptra_attrID  FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_yptra_dictID  FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_yptra_uq_dict FOREIGN KEY     ( trcID, dictionaryID ) REFERENCES yp.technology_reference_card ( trcID, dictionaryID ),
    CONSTRAINT __fk_yptra_uq_att  FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __yptra_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(trcID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS yp.technology_reference_card_unique_attribute_values (
    trcID                         uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_yptrq_serID   FOREIGN KEY     ( trcID ) REFERENCES yp.technology_reference_card ( trcID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_yptrq_attrID  FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_yptrq_dictID  FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_yptrq_uq_dict FOREIGN KEY     ( trcID, dictionaryID ) REFERENCES yp.technology_reference_card ( trcID, dictionaryID ),
    CONSTRAINT __fk_yptrq_uq_att  FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __yptrq_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(trcID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __yptrq_temp_uniq  EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
