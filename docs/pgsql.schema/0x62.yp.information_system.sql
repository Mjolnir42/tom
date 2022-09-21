--
--
-- YP SCHEMA
CREATE TABLE IF NOT EXISTS yp.information_system (
    isID                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_ixis          PRIMARY KEY     ( isID ),
    CONSTRAINT __fk_ixis_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixis_fk_origin   UNIQUE          ( isID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS yp.information_system_standard_attribute_values (
    isID                          uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_ixisa_isID    FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixisa_attrID  FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixisa_dictID  FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixisa_uq_dct  FOREIGN KEY     ( isID, dictionaryID ) REFERENCES yp.information_system ( isID, dictionaryID ),
    CONSTRAINT __fk_ixisa_uq_att  FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixisa_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(isID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS yp.information_system_unique_attribute_values (
    isID                          uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_ixisq_isID    FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixisq_attrID  FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixisq_dictID  FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixisq_uq_dct  FOREIGN KEY     ( isID, dictionaryID ) REFERENCES yp.information_system ( isID, dictionaryID ),
    CONSTRAINT __fk_ixisq_uq_att  FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixisq_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(isID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __ixisq_temp_uniq  EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
