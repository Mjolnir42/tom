---
---
---
CREATE TABLE IF NOT EXISTS abstract.data_parent (
    dataID                        uuid            NOT NULL,
    dataDictionaryID              uuid            NOT NULL,
    bpoID                         uuid            NULL,
    bpoDictionaryID               uuid            NULL,
    moduleID                      uuid            NULL,
    moduleDictionaryID            uuid            NULL,
    artifactID                    uuid            NULL,
    artifactDictionaryID          uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_adp_dataID    FOREIGN KEY     ( dataID, dataDictionaryID ) REFERENCES abstract.data ( dataID, dictionaryID ),
    CONSTRAINT __fk_adp_bpoID     FOREIGN KEY     ( bpoID, bpoDictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
    CONSTRAINT __fk_adp_modID     FOREIGN KEY     ( moduleID, moduleDictionaryID ) REFERENCES abstract.module ( moduleID, dictionaryID ),
    CONSTRAINT __fk_adp_artID     FOREIGN KEY     ( artifactID, artifactDictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __adp_nonnull      CHECK           (   ((bpoID IS NOT NULL) AND (moduleID IS     NULL) AND (artifactID IS     NULL))
                                                   OR ((bpoID IS     NULL) AND (moduleID IS NOT NULL) AND (artifactID IS     NULL))
                                                   OR ((bpoID IS     NULL) AND (moduleID IS     NULL) AND (artifactID IS NOT NULL))),
    CONSTRAINT __adp_null_bpo     CHECK           (   ((bpoID IS NOT NULL) AND (bpoDictionaryID IS NOT NULL))
                                                   OR ((bpoID IS     NULL) AND (bpoDictionaryID IS     NULL))),
    CONSTRAINT __adp_null_dpl     CHECK           (   ((moduleID IS NOT NULL) AND (moduleDictionaryID IS NOT NULL))
                                                   OR ((moduleID IS     NULL) AND (moduleDictionaryID IS     NULL))),
    CONSTRAINT __adp_null_ins     CHECK           (   ((artifactID IS NOT NULL) AND (artifactDictionaryID IS NOT NULL))
                                                   OR ((artifactID IS     NULL) AND (artifactDictionaryID IS     NULL))),
    CONSTRAINT __adp_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(dataID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __adp_temp_bpo     EXCLUDE         USING gist (public.uuid_to_bytea(dataID) WITH =,
                                                              public.uuid_to_bytea(bpoID) WITH =,
                                                              validity WITH &&) WHERE (bpoID IS NOT NULL),
    CONSTRAINT __adp_temp_module  EXCLUDE         USING gist (public.uuid_to_bytea(dataID) WITH =,
                                                              public.uuid_to_bytea(moduleID) WITH =,
                                                              validity WITH &&) WHERE (moduleID IS NOT NULL),
    CONSTRAINT __adp_temp_ins     EXCLUDE         USING gist (public.uuid_to_bytea(dataID) WITH =,
                                                              public.uuid_to_bytea(artifactID) WITH =,
                                                              validity WITH &&) WHERE (artifactID IS NOT NULL)
);
CREATE TABLE IF NOT EXISTS abstract.service_parent (
    serviceID                     uuid            NOT NULL,
    serviceDictionaryID           uuid            NOT NULL,
    bpoID                         uuid            NULL,
    bpoDictionaryID               uuid            NULL,
    moduleID                      uuid            NULL,
    moduleDictionaryID            uuid            NULL,
    artifactID                    uuid            NULL,
    artifactDictionaryID          uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange((NOW() AT TIME ZONE 'utc'), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT __fk_asp_serviceID FOREIGN KEY     ( serviceID, serviceDictionaryID ) REFERENCES abstract.service ( serviceID, dictionaryID ),
    CONSTRAINT __fk_asp_bpoID     FOREIGN KEY     ( bpoID, bpoDictionaryID ) REFERENCES abstract.blueprint ( bpoID, dictionaryID ),
    CONSTRAINT __fk_asp_modID     FOREIGN KEY     ( moduleID, moduleDictionaryID ) REFERENCES abstract.module ( moduleID, dictionaryID ),
    CONSTRAINT __fk_asp_artID     FOREIGN KEY     ( artifactID, artifactDictionaryID ) REFERENCES abstract.artifact ( artifactID, dictionaryID ),
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __asp_nonnull      CHECK           (   ((bpoID IS NOT NULL) AND (moduleID IS     NULL) AND (artifactID IS     NULL))
                                                   OR ((bpoID IS     NULL) AND (moduleID IS NOT NULL) AND (artifactID IS     NULL))
                                                   OR ((bpoID IS     NULL) AND (moduleID IS     NULL) AND (artifactID IS NOT NULL))),
    CONSTRAINT __asp_null_bpo     CHECK           (   ((bpoID IS NOT NULL) AND (bpoDictionaryID IS NOT NULL))
                                                   OR ((bpoID IS     NULL) AND (bpoDictionaryID IS     NULL))),
    CONSTRAINT __asp_null_dpl     CHECK           (   ((moduleID IS NOT NULL) AND (moduleDictionaryID IS NOT NULL))
                                                   OR ((moduleID IS     NULL) AND (moduleDictionaryID IS     NULL))),
    CONSTRAINT __asp_null_ins     CHECK           (   ((artifactID IS NOT NULL) AND (artifactDictionaryID IS NOT NULL))
                                                   OR ((artifactID IS     NULL) AND (artifactDictionaryID IS     NULL))),
    CONSTRAINT __asp_temporal     EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                              validity WITH &&),
    CONSTRAINT __asp_temp_bpo     EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                              public.uuid_to_bytea(bpoID) WITH =,
                                                              validity WITH &&) WHERE (bpoID IS NOT NULL),
    CONSTRAINT __asp_temp_module  EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                              public.uuid_to_bytea(moduleID) WITH =,
                                                              validity WITH &&) WHERE (moduleID IS NOT NULL),
    CONSTRAINT __asp_temp_ins     EXCLUDE         USING gist (public.uuid_to_bytea(serviceID) WITH =,
                                                              public.uuid_to_bytea(artifactID) WITH =,
                                                              validity WITH &&) WHERE (artifactID IS NOT NULL)
);
