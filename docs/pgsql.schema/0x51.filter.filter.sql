---
---
--- FILTER SCHEMA
CREATE TABLE IF NOT EXISTS filter.filter (
    filterID                      uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    dictionaryID                  uuid        NOT NULL,
    CONSTRAINT __pk_ff            PRIMARY KEY ( filterID ),
    CONSTRAINT __fk_ff_dictID     FOREIGN KEY ( dictionaryID ) REFERENCES meta.dictionary ( dictionaryID ) ON DELETE RESTRICT,
    CONSTRAINT __ff_fk_origin     UNIQUE      ( filterID, dictionaryID )
);
CREATE TABLE IF NOT EXISTS filter.name (
    filterID                      uuid        NOT NULL,
    dictionaryID                  uuid        NOT NULL,
    name                          text        NOT NULL,
    cardinality                   flt_card    NOT NULL DEFAULT 'one'::flt_card,
    aggregation                   flt_aggr    NOT NULL DEFAULT 'max'::flt_aggr,
    CONSTRAINT __fk_ffn_origin    FOREIGN KEY ( filterID, dictionaryID ) REFERENCES filter.filter ( filterID, dictionaryID ) DEFERRABLE,
    CONSTRAINT __ffn_fk_card      UNIQUE      ( filterID, cardinality ),
    CONSTRAINT __ffn_uniq_name    UNIQUE      ( dictionaryID, name )
);
CREATE TABLE IF NOT EXISTS filter.value (
    filterValueID                 uuid        NOT NULL DEFAULT public.gen_random_uuid(),
    filterID                      uuid        NOT NULL,
    value                         text        NOT NULL,
    valueOrder                    smallint    NOT NULL DEFAULT 0,
    CONSTRAINT __pk_ffv           PRIMARY KEY ( filterValueID ),
    CONSTRAINT __fk_ffv_filterID  FOREIGN KEY ( filterID ) REFERENCES filter.filter ( filterID ) DEFERRABLE,
    CONSTRAINT __ffv_uniq_value   UNIQUE      ( filterID, value ),
    CONSTRAINT __ffv_fk_origin    UNIQUE      ( filterValueID, filterID )
);
CREATE TABLE IF NOT EXISTS filter.assignable_entity (
    filterID                      uuid        NOT NULL,
    entity                        flt_ntt     NOT NULL,
    CONSTRAINT __fk_ffae_filterID FOREIGN KEY ( filterID ) REFERENCES filter.filter ( filterID ) DEFERRABLE,
    CONSTRAINT __ffae_fk_origin   UNIQUE      ( filterID, entity )
);
