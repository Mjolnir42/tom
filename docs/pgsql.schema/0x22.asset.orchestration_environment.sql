--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.orchestration_environment (
    orchID                        uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __pk_aoe           PRIMARY KEY     ( orchID ),
    CONSTRAINT __fk_aoe_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aoe_fk_origin    UNIQUE          ( orchID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS asset.orchestration_environment_linking (
    orchestrationLinkID           uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    orchID_A                      uuid            NOT NULL,
    dictionaryID_A                uuid            NOT NULL,
    orchID_B                      uuid            NOT NULL,
    dictionaryID_B                uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __pk_aoel          PRIMARY KEY     ( orchestrationLinkID ),
    CONSTRAINT __fk_aoel_rteA     FOREIGN KEY     ( orchID_A, dictionaryID_A ) REFERENCES asset.orchestration_environment ( orchID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoel_rteB     FOREIGN KEY     ( orchID_B, dictionaryID_B ) REFERENCES asset.orchestration_environment ( orchID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aoel_diff_rte    CHECK           ( orchID_A != orchID_B ),
    CONSTRAINT __aoel_uniq_link   UNIQUE          ( orchID_A, orchID_B ),
    CONSTRAINT __aeol_ordered     CHECK           (public.uuid_to_bytea(orchID_A) > public.uuid_to_bytea(orchID_B))
);
CREATE TABLE IF NOT EXISTS asset.orchestration_environment_standard_attribute_values (
    orchID                        uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __fk_aoea_orchID   FOREIGN KEY     ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoea_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoea_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoea_uq_dict  FOREIGN KEY     ( orchID, dictionaryID ) REFERENCES asset.orchestration_environment ( orchID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoea_uq_attr  FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aoea_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(orchID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.orchestration_environment_unique_attribute_values (
    orchID                        uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT __fk_aoeq_orchID   FOREIGN KEY     ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoeq_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoeq_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoeq_uq_dic   FOREIGN KEY     ( orchID, dictionaryID ) REFERENCES asset.orchestration_environment ( orchID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_aoeq_uq_att   FOREIGN KEY     ( attributeID, dictionaryID ) REFERENCES meta.unique_attribute ( attributeID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __aoeq_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(orchID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __aoeq_temp_uniq   EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
