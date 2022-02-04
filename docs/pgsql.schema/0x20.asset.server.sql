--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.server (
    serverID                      uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_as            PRIMARY KEY     ( serverID ),
    CONSTRAINT __fk_as_dictID     FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __as_fk_origin     UNIQUE          ( serverID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS asset.server_linking (
    serverLinkID                  uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    serverID_A                    uuid            NOT NULL,
    dictionaryID_A                uuid            NOT NULL,
    serverID_B                    uuid            NOT NULL,
    dictionaryID_B                uuid            NOT NULL,
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_asl           PRIMARY KEY     ( serverLinkID ),
    CONSTRAINT __fk_asl_srvA      FOREIGN KEY     ( serverID_A, dictionaryID_A ) REFERENCES asset.server ( serverID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asl_srvB      FOREIGN KEY     ( serverID_B, dictionaryID_B ) REFERENCES asset.server ( serverID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __asl_diff_servID  CHECK           ( serverID_A != serverID_B ),
    CONSTRAINT __asl_uniq_link    UNIQUE          ( serverID_A, serverID_B ),
    CONSTRAINT __asl_ordered      CHECK           ( public.uuid_to_bytea(serverID_A) > public.uuid_to_bytea(serverID_B))
);
CREATE TABLE IF NOT EXISTS asset.server_standard_attribute_values (
    serverID                      uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_asa_servID    FOREIGN KEY     ( serverID ) REFERENCES asset.server ( serverID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asa_attrID    FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asa_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asa_uq_dict   FOREIGN KEY     ( serverID, dictionaryID ) REFERENCES asset.server ( serverID, dictionaryID ),
    CONSTRAINT __fk_asa_uq_attr   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __asa_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(serverID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.server_unique_attribute_values (
    serverID                      uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_asq_servID    FOREIGN KEY     ( serverID ) REFERENCES asset.server ( serverID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asq_attrID    FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asq_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asq_uq_dict   FOREIGN KEY     ( serverID, dictionaryID ) REFERENCES asset.server ( serverID, dictionaryID ),
    CONSTRAINT __fk_asq_uq_attr   FOREIGN KEY     ( attributeID, dictionaryID ) REFERENCES meta.unique_attribute ( attributeID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __asq_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(serverID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&),
    CONSTRAINT __asq_temp_uniq    EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                          public.uuid_to_bytea(dictionaryID) WITH =,
                                                          value WITH =,
                                                          validity WITH &&)
);
