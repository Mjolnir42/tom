--
--
-- YP SCHEMA
CREATE TABLE IF NOT EXISTS yp.corporate_domain_parent (
    corporateID                   uuid            NOT NULL,
    domainID                      uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
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
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
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
    serID                         uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_tomd_isID     FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_tomd_serID    FOREIGN KEY     ( serID ) REFERENCES yp.ypservice ( serID ) ON DELETE RESTRICT,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __tois_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(serID) WITH =,
                                                              validity WITH &&)
);
