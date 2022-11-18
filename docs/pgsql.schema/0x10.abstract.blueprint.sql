--
--
-- abstract SCHEMA
CREATE TABLE IF NOT EXISTS abstract.blueprint (
    bpoID                         uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_abpo          PRIMARY KEY     ( bpoID ),
    CONSTRAINT __fk_abpo_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __abpo_fk_origin   UNIQUE          ( bpoID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS abstract.blueprint_standard_attribute_values (
    bpoID                         uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_abpoav_bpoID  FOREIGN KEY     ( bpoID ) REFERENCES abstract.blueprint ( bpoID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_abpoav_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_abpoav_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_abpoav_uq_dct FOREIGN KEY     ( bpoID, dictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
    CONSTRAINT __fk_abpoav_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __abpoav_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(bpoID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS abstract.blueprint_unique_attribute_values (
    bpoID                         uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_abpoqv_bpoID  FOREIGN KEY     ( bpoID ) REFERENCES abstract.blueprint ( bpoID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_abpoqv_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_abpoqv_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_abpoqv_uq_dct FOREIGN KEY     ( bpoID, dictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
    CONSTRAINT __fk_abpoqv_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __abpoqv_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(bpoID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __abpoqv_temp_uniq EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
