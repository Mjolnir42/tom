---
---
---
CREATE TABLE IF NOT EXISTS abstract.blueprint_mapping (
    bpoID                         uuid            NOT NULL,
    bpoDictionaryID               uuid            NOT NULL,
    moduleID                      uuid            NOT NULL,
    moduleDictionaryID            uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_abm_bpoID     FOREIGN KEY     ( bpoID, bpoDictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
    CONSTRAINT __fk_abm_modID     FOREIGN KEY     ( moduleID, moduleDictionaryID ) REFERENCES abstract.module ( moduleID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __abm_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(bpoID) WITH =,
                                                              public.uuid_to_bytea(moduleID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS abstract.module_mapping (
    moduleID                      uuid            NOT NULL,
    moduleDictionaryID            uuid            NOT NULL,
    artifactID                    uuid            NOT NULL,
    artifactDictionaryID          uuid            NOT NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_amm_modID     FOREIGN KEY     ( moduleID, moduleDictionaryID ) REFERENCES abstract.module ( moduleID, dictionaryID ),
    CONSTRAINT __fk_amm_artID     FOREIGN KEY     ( artifactID, artifactDictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __amm_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(moduleID) WITH =,
                                                              public.uuid_to_bytea(artifactID) WITH =,
                                                              validity WITH &&)
);
