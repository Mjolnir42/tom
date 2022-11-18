--
--
-- iX SCHEMA
CREATE TABLE IF NOT EXISTS ix.consumer_product (
    productID                     uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __pk_ixp           PRIMARY KEY     ( productID ),
    CONSTRAINT __fk_ixp_dictID    FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixp_fk_origin    UNIQUE          ( productID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS ix.consumer_product_standard_attribute_values (
    productID                     uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_ixpa_prodID   FOREIGN KEY     ( productID ) REFERENCES ix.consumer_product ( productID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixpa_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.standard_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixpa_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixpa_uq_dct   FOREIGN KEY     ( productID, dictionaryID ) REFERENCES ix.consumer_product ( productID, dictionaryID ),
    CONSTRAINT __fk_ixpa_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.standard_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixpa_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(productID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS ix.consumer_product_unique_attribute_values (
    productID                     uuid            NOT NULL,
    attributeID                   uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    value                         text            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_ixpq_prodID   FOREIGN KEY     ( productID ) REFERENCES ix.consumer_product ( productID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixpq_attrID   FOREIGN KEY     ( attributeID ) REFERENCES meta.unique_attribute ( attributeID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixpq_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixpq_uq_dct   FOREIGN KEY     ( productID, dictionaryID ) REFERENCES ix.consumer_product ( productID, dictionaryID ),
    CONSTRAINT __fk_ixpq_uq_att   FOREIGN KEY     ( dictionaryID, attributeID ) REFERENCES meta.unique_attribute ( dictionaryID, attributeID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixpq_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(productID) WITH =,
                                                              public.uuid_to_bytea(attributeID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __ixpq_temp_uniq   EXCLUDE         USING gist (public.uuid_to_bytea(attributeID) WITH =,
                                                              public.uuid_to_bytea(dictionaryID) WITH =,
                                                              value WITH =,
                                                              validity WITH &&)
);
