--
--
-- ASSET INVENTORY DATA
CREATE TABLE IF NOT EXISTS asset.socket (
    socketID                      uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid        NOT NULL,
    CONSTRAINT __pk_ass           PRIMARY KEY ( socketID ),
    CONSTRAINT __fk_ass__dictID   FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __ass_fk_origin    UNIQUE      ( socketID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS asset.socket_linking (
    socketLinkID                  uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    socketID_A                    uuid        NOT NULL,
    dictionaryID_A                uuid        NOT NULL,
    socketID_B                    uuid        NOT NULL,
    dictionaryID_B                uuid        NOT NULL,
    CONSTRAINT __pk_assl          PRIMARY KEY ( socketLinkID ),
    CONSTRAINT __fk_assl_sockA    FOREIGN KEY ( socketID_A, dictionaryID_A ) REFERENCES asset.socket ( socketID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_assl_sockB    FOREIGN KEY ( socketID_B, dictionaryID_B ) REFERENCES asset.socket ( socketID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __assl_diff_sock   CHECK       ( socketID_A != socketID_B ),
    CONSTRAINT __assl_uniq_link   UNIQUE      ( socketID_A, socketID_B ),
    CONSTRAINT __assl_ordered     CHECK       ( public.uuid_to_bytea(socketID_A) > public.uuid_to_bytea(socketID_B))
);
CREATE TABLE IF NOT EXISTS asset.socket_standard_attribute_values (
    socketID                      uuid        NOT NULL,
    attributeID                   uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    value                         text        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_asssav_sockID FOREIGN KEY ( socketID ) REFERENCES asset.socket ( socketID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asssav_attrID FOREIGN KEY ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asssav_dictID FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asssav_uq_dic FOREIGN KEY ( socketID, dictionaryID ) REFERENCES asset.socket ( socketID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_asssav_uq_att FOREIGN KEY ( attributeID, dictionaryID ) REFERENCES meta.standard_attribute ( attributeID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __asssav_temporal  EXCLUDE     USING gist (public.uuid_to_bytea(socketID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&)
);
CREATE TABLE IF NOT EXISTS asset.socket_unique_attribute_values (
    socketID                      uuid        NOT NULL,
    attributeID                   uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    value                         text        NOT NULL,
    validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    CONSTRAINT __fk_assuav_sockID FOREIGN KEY ( socketID ) REFERENCES asset.socket ( socketID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_assuav_attrID FOREIGN KEY ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_assuav_dictID FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_assuav_uq_dic FOREIGN KEY ( socketID, dictionaryID ) REFERENCES asset.socket ( socketID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_assuav_uq_att FOREIGN KEY ( attributeID, dictionaryID ) REFERENCES meta.unique_attribute ( attributeID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __assuav_temporal  EXCLUDE     USING gist (public.uuid_to_bytea(socketID) WITH =,
                                                          public.uuid_to_bytea(attributeID) WITH =,
                                                          validity WITH &&),
    CONSTRAINT __assuav_temp_uniq EXCLUDE     USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                          public.uuid_to_bytea(dictionaryID) WITH =,
                                                          value WITH =,
                                                          validity WITH &&)
);
