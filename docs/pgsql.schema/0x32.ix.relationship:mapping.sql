--
--
-- iX SCHEMA
CREATE TABLE IF NOT EXISTS ix.top_level_service_mapping (
    tlsID                         uuid            NOT NULL,
    isID                          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_ixmtls_tlsID  FOREIGN KEY     ( tlsID ) REFERENCES ix.top_level_service ( tlsID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixmtls_isID   FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixmtls_temporal  EXCLUDE         USING gist (public.uuid_to_bytea(tlsID) WITH =,
                                                              public.uuid_to_bytea(isID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS ix.product_mapping (
    productID                     uuid            NOT NULL,
    isID                          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_ixmp_prodID   FOREIGN KEY     ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixmp_isID     FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixmp_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(productID) WITH =,
                                                              public.uuid_to_bytea(isID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS ix.endpoint_mapping (
    endpointID                    uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    componentID                   uuid            NULL,
    componentDictionaryID         uuid            NULL,
    groupID                       uuid            NULL,
    groupDictionaryID             uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_ixme_endpID   FOREIGN KEY     ( endpointID ) REFERENCES ix.endpoint ( endpointID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_dictID   FOREIGN KEY     ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_uq_endp  FOREIGN KEY     ( endpointID, dictionaryID ) REFERENCES ix.endpoint ( endpointID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_cmpID    FOREIGN KEY     ( componentID ) REFERENCES ix.functional_component ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_cDictID  FOREIGN KEY     ( componentDictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_uq_comp  FOREIGN KEY     ( componentID, componentDictionaryID ) REFERENCES ix.functional_component ( componentID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_grpID    FOREIGN KEY     ( groupID ) REFERENCES ix.deployment_group ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_gDictID  FOREIGN KEY     ( groupDictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixme_uq_group FOREIGN KEY     ( groupID, groupDictionaryID ) REFERENCES ix.deployment_group ( groupID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixme_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(endpointID) WITH =,
                                                              public.uuid_to_bytea(componentID) WITH =,
                                                              public.uuid_to_bytea(groupID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __ixme_uniq_map    CHECK           (   ((componentID IS NOT NULL) AND (componentDictionaryID IS NOT NULL) AND (groupID IS     NULL) AND (groupDictionaryID IS     NULL))
                                                   OR ((componentID IS     NULL) AND (componentDictionaryID IS     NULL) AND (groupID IS NOT NULL) AND (groupDictionaryID IS NOT NULL)))
);
CREATE TABLE IF NOT EXISTS ix.deployment_group_mapping (
    groupID                       uuid            NOT NULL,
    techsrvID                     uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_ixmdg_grpID   FOREIGN KEY     ( groupID ) REFERENCES ix.deployment_group ( groupID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixmdg_techID  FOREIGN KEY     ( techsrvID ) REFERENCES ix.technical_service ( techsrvID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixmdg_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(groupID) WITH =,
                                                              public.uuid_to_bytea(techsrvID) WITH =,
                                                              validity WITH &&)
);
