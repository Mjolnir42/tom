--
--
-- YP SCHEMA
CREATE TABLE IF NOT EXISTS yp.corporate_domain_parent (
    corporateID                   uuid            NOT NULL,
    domainID                      uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_tomcd_corpID  FOREIGN KEY     ( corporateID ) REFERENCES yp.corporate_domain ( corporateID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_tomcd_domID   FOREIGN KEY     ( domainID ) REFERENCES yp.domain ( domainID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __tomcd_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(domainID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS yp.domain_parent (
    domainID                      uuid            NOT NULL,
    isID                          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_tomd_domID    FOREIGN KEY     ( domainID ) REFERENCES yp.domain ( domainID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_tomd_isID     FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __tomd_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(isID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS yp.information_system_parent (
    isID                          uuid            NOT NULL,
    componentID                   uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_ixmis_isID    FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ixmis_compID  FOREIGN KEY     ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ixmis_temporal   EXCLUDE         USING gist (public.uuid_to_bytea(componentID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS yp.service_parent (
    serviceID                     uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    parentInformationSystemID     uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_ypsp_srvID    FOREIGN KEY     ( serviceID, dictionaryID ) REFERENCES yp.service ( serviceID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_ypsp_isID     FOREIGN KEY     ( parentInformationSystemID, dictionaryID ) REFERENCES yp.information_system ( isID, dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __ypsp_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                              validity WITH &&)
);
