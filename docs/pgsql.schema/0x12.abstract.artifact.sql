--
--
-- abstract SCHEMA
CREATE TABLE IF NOT EXISTS abstract.artifact (
    artifactID                    uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_aart          PRIMARY KEY     ( artifactID ),
    CONSTRAINT __fk_aart_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aart_fk_origin   UNIQUE          ( artifactID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS abstract.artifact_standard_attribute_values (
    artifactID                    uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_aartav_artfID FOREIGN KEY     ( artifactID ) REFERENCES abstract.artifact ( artifactID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aartav_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aartav_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aartav_uq_dct FOREIGN KEY     ( artifactID, dictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
    CONSTRAINT __fk_aartav_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aartav_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(artifactID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS abstract.artifact_unique_attribute_values (
    artifactID                    uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_aartqv_artfID FOREIGN KEY     ( artifactID ) REFERENCES abstract.artifact ( artifactID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aartqv_attrID FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aartqv_dictID FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aartqv_uq_dct FOREIGN KEY     ( artifactID, dictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
    CONSTRAINT __fk_aartqv_uq_att FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aartqv_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(artifactID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __aartqv_temp_uniq EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
