--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.runtime_environment (
    rteID                         uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_arte          PRIMARY KEY     ( rteID ),
    CONSTRAINT __fk_arte_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __arte_fk_origin   UNIQUE          ( rteID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS asset.runtime_environment_linking (
    runtimeLinkID                 uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    rteID_A                       uuid            NOT NULL,
    dictionaryID_A                uuid            NOT NULL,
    rteID_B                       uuid            NOT NULL,
    dictionaryID_B                uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_artel         PRIMARY KEY     ( runtimeLinkID ),
    CONSTRAINT __fk_artel_rteA    FOREIGN KEY     ( rteID_A, dictionaryID_A ) REFERENCES asset.runtime_environment ( rteID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artel_rteB    FOREIGN KEY     ( rteID_B, dictionaryID_B ) REFERENCES asset.runtime_environment ( rteID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __artel_diff_rte   CHECK           ( rteID_A != rteID_B ),
    CONSTRAINT __artel_uniq_link  UNIQUE          ( rteID_A, rteID_B ),
    CONSTRAINT __artel_ordered    CHECK           ( public.uuid_to_bytea(rteID_A) > public.uuid_to_bytea(rteID_B))
);
CREATE TABLE IF NOT EXISTS asset.runtime_environment_standard_attribute_values (
    rteID                         uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_artea_rteID   FOREIGN KEY     ( rteID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artea_attrID  FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artea_dictID  FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artea_uq_dic  FOREIGN KEY     ( rteID, dictionaryID ) REFERENCES asset.runtime_environment ( rteID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_artea_uq_att  FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __artea_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(rteID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.runtime_environment_unique_attribute_values (
    rteID                         uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_arteq_rteID   FOREIGN KEY     ( rteID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_arteq_attrID  FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_arteq_dictID  FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_arteq_uq_dic  FOREIGN KEY     ( rteID, dictionaryID ) REFERENCES asset.runtime_environment ( rteID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_arteq_uq_att  FOREIGN KEY     ( attributeID, dictionaryID ) REFERENCES meta.unique_attribute ( attributeID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __arteq_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(rteID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __arteq_temp_uniq  EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
