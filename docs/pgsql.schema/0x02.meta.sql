--
--
-- DATABASE META DATA --
CREATE TABLE IF NOT EXISTS meta.dictionary (
    dictionaryID                  uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    -- the current name of the dictionary is duplicated here for bootstrap
    -- reasons. to ensure unique dictionary names would otherwise
    -- require a dictionary of dictionaries to enforce unique naming.
    -- this would mean turtles all the way down.
    name                          text            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __pk_mdict         PRIMARY KEY     ( dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __uniq_dictionary  UNIQUE          ( name ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
);
CREATE TABLE IF NOT EXISTS meta.attribute (
    dictionaryID                  uuid            NOT NULL,
    attribute                     text            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __uniq_attr_name   UNIQUE          ( dictionaryID, attribute ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
);
CREATE TABLE IF NOT EXISTS meta.standard_attribute (
    attributeID                   uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    attribute                     text            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __pk_msa           PRIMARY KEY     ( attributeID ),
    CONSTRAINT __fk_msa_attr      FOREIGN KEY     ( dictionaryID, attribute ) REFERENCES meta.attribute ( dictionaryID, attribute )
                                                  ON DELETE CASCADE,
    CONSTRAINT __fk_msa_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __uniq_attribute   UNIQUE          ( dictionaryID, attribute ),
    CONSTRAINT __msa_fk_origin    UNIQUE          ( dictionaryID, attributeID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
);
CREATE TABLE IF NOT EXISTS meta.unique_attribute (
    attributeID                   uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    attribute                     text            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __pk_msqa          PRIMARY KEY     ( attributeID ),
    CONSTRAINT __fk_msqa_attr     FOREIGN KEY     ( dictionaryID, attribute ) REFERENCES meta.attribute ( dictionaryID, attribute )
                                                  ON DELETE CASCADE,
    CONSTRAINT __fk_msqa_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __uniq_unique_attr UNIQUE          ( dictionaryID, attribute ),
    CONSTRAINT __msqa_fk_origin   UNIQUE          ( dictionaryID, attributeID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' )
);
CREATE TABLE IF NOT EXISTS meta.dictionary_standard_attribute_values (
    dictionaryID                  uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_mda_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ),
    CONSTRAINT __fk_mda_attrID    FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ),
    CONSTRAINT __fk_mda_uq_attr   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __mda_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(dictionaryID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS meta.dictionary_unique_attribute_values (
    dictionaryID                  uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_mdq_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ),
    CONSTRAINT __fk_mdq_attrID    FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ),
    CONSTRAINT __fk_mdq_uq_attr   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __mdq_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(dictionaryID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __mdq_temp_uniq    EXCLUDE         USING gist (public.uuid_to_bytea(dictionaryID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
