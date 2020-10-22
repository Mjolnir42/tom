--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.container (
    containerID                   uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid        NOT NULL,
    CONSTRAINT __pk_asc           PRIMARY KEY ( containerID ),
    CONSTRAINT __fk_asc__dictID   FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __asc_fk_origin    UNIQUE      ( containerID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS asset.container_linking (
    containerLinkID               uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    containerID_A                 uuid        NOT NULL,
    dictionaryID_A                uuid        NOT NULL,
    containerID_B                 uuid        NOT NULL,
    dictionaryID_B                uuid        NOT NULL,
    CONSTRAINT __pk_ascl          PRIMARY KEY ( containerLinkID ),
    CONSTRAINT __fk_ascl_sockA    FOREIGN KEY ( containerID_A, dictionaryID_A ) REFERENCES asset.container ( containerID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascl_sockB    FOREIGN KEY ( containerID_B, dictionaryID_B ) REFERENCES asset.container ( containerID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __ascl_diff_sock   CHECK       ( containerID_A != containerID_B ),
    CONSTRAINT __ascl_uniq_link   UNIQUE      ( containerID_A, containerID_B ),
    CONSTRAINT __ascl_ordered     CHECK       ( public.uuid_to_bytea(containerID_A) > public.uuid_to_bytea(containerID_B))
);
CREATE TABLE IF NOT EXISTS asset.container_standard_attribute_values (
    containerID                   uuid        NOT NULL,
    attributeID                   uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    value                         text        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_ascsav_sockID FOREIGN KEY ( containerID ) REFERENCES asset.container ( containerID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascsav_attrID FOREIGN KEY ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascsav_dictID FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascsav_uq_dic FOREIGN KEY ( containerID, dictionaryID ) REFERENCES asset.container ( containerID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascsav_uq_att FOREIGN KEY ( attributeID, dictionaryID ) REFERENCES meta.standard_attribute ( attributeID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __ascsav_temporal  EXCLUDE     USING gist (public.uuid_to_bytea(containerID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.container_unique_attribute_values (
    containerID                   uuid        NOT NULL,
    attributeID                   uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    value                         text        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_ascuav_sockID FOREIGN KEY ( containerID ) REFERENCES asset.container ( containerID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascuav_attrID FOREIGN KEY ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascuav_dictID FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascuav_uq_dic FOREIGN KEY ( containerID, dictionaryID ) REFERENCES asset.container ( containerID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ascuav_uq_att FOREIGN KEY ( attributeID, dictionaryID ) REFERENCES meta.unique_attribute ( attributeID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __ascuav_temporal  EXCLUDE     USING gist (public.uuid_to_bytea(containerID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&),
    CONSTRAINT __ascuav_temp_uniq EXCLUDE     USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                          public.uuid_to_bytea(dictionaryID) WITH =,
                                                          value WITH =,
                                                          validity WITH &&)
);
