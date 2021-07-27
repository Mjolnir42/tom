---
---
--- FILTER SCHEMA
CREATE TABLE IF NOT EXISTS filter.value_assignment__one (
    filterValueID                 uuid            NOT NULL,
    filterID                      uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    cardinality                   flt_card        NOT NULL,
    entity                        flt_ntt         NOT NULL,
    tlsID                         uuid            NULL,
    productID                     uuid            NULL,
    isID                          uuid            NULL,
    componentID                   uuid            NULL,
    groupID                       uuid            NULL,
    orchID                        uuid            NULL,
    rteID                         uuid            NULL,
    serverID                      uuid            NULL,
    endpointID                    uuid            NULL,
    containerID                   uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_fvao_value    FOREIGN KEY     ( filterValueID, filterID ) REFERENCES filter.value ( filterValueID, filterID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_filter   FOREIGN KEY     ( filterID, dictionaryID ) REFERENCES filter.filter ( filterID, dictionaryID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_card     FOREIGN KEY     ( filterID, cardinality ) REFERENCES filter.name ( filterID, cardinality ) DEFERRABLE,
    CONSTRAINT __fk_fvao_assign   FOREIGN KEY     ( filterID, entity ) REFERENCES filter.assignable_entity ( filterID, entity ) DEFERRABLE,
    CONSTRAINT __fk_fvao_tlsID    FOREIGN KEY     ( tlsID ) REFERENCES ix.top_level_service ( tlsID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_prodID   FOREIGN KEY     ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_isID     FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_compID   FOREIGN KEY     ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_groupID  FOREIGN KEY     ( groupID ) REFERENCES ix.deployment_group ( groupID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvao_orchID   FOREIGN KEY     ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_rteID    FOREIGN KEY     ( rteID ) REFERENCES asset.runtime_environment ( rteID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_serverID FOREIGN KEY     ( serverID ) REFERENCES asset.server ( serverID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_endpID   FOREIGN KEY     ( endpointID ) REFERENCES ix.endpoint ( endpointID ) DEFERRABLE,
    CONSTRAINT __fk_fvao_contID   FOREIGN KEY     ( containerID ) REFERENCES asset.container ( containerID ) DEFERRABLE,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __fvao_cardinality CHECK           ( cardinality = 'one'::flt_card ),
    CONSTRAINT __fvao_uniq_object CHECK           (   ((entity='top_level_service')         AND (tlsID       IS NOT NULL))
                                                   OR ((entity='product')                   AND (productID   IS NOT NULL))
                                                   OR ((entity='information_system')        AND (isID        IS NOT NULL))
                                                   OR ((entity='functional_component')      AND (componentID IS NOT NULL))
                                                   OR ((entity='deployment_group')          AND (groupID     IS NOT NULL))
                                                   OR ((entity='orchestration_environment') AND (orchID      IS NOT NULL))
                                                   OR ((entity='runtime_environment')       AND (rteID       IS NOT NULL))
                                                   OR ((entity='server')                    AND (serverID    IS NOT NULL))
                                                   OR ((entity='endpoint')                  AND (endpointID  IS NOT NULL))
                                                   OR ((entity='container')                 AND (containerID IS NOT NULL))),
    CONSTRAINT __fvao_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(filterID) WITH =,
                                                              public.uuid_to_bytea(tlsID) WITH =,
                                                              public.uuid_to_bytea(productID) WITH =,
                                                              public.uuid_to_bytea(isID) WITH =,
                                                              public.uuid_to_bytea(componentID) WITH =,
                                                              public.uuid_to_bytea(groupID) WITH =,
                                                              public.uuid_to_bytea(orchID) WITH =,
                                                              public.uuid_to_bytea(rteID) WITH =,
                                                              public.uuid_to_bytea(serverID) WITH =,
                                                              public.uuid_to_bytea(endpointID) WITH =,
                                                              public.uuid_to_bytea(containerID) WITH =,
                                                              validity WITH &&)
);
CREATE TABLE IF NOT EXISTS filter.value_assignment__many (
    filterValueID                 uuid            NOT NULL,
    filterID                      uuid            NOT NULL,
    dictionaryID                  uuid            NOT NULL,
    cardinality                   flt_card        NOT NULL,
    entity                        flt_ntt         NOT NULL,
    tlsID                         uuid            NULL,
    productID                     uuid            NULL,
    isID                          uuid            NULL,
    componentID                   uuid            NULL,
    groupID                       uuid            NULL,
    orchID                        uuid            NULL,
    rteID                         uuid            NULL,
    serverID                      uuid            NULL,
    endpointID                    uuid            NULL,
    containerID                   uuid            NULL,
    validity                      tstzrange       NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
    createdBy                     uuid            NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT __fk_fvam_value    FOREIGN KEY     ( filterValueID, filterID ) REFERENCES filter.value ( filterValueID, filterID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_filter   FOREIGN KEY     ( filterID, dictionaryID ) REFERENCES filter.filter ( filterID, dictionaryID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_card     FOREIGN KEY     ( filterID, cardinality ) REFERENCES filter.name ( filterID, cardinality ) DEFERRABLE,
    CONSTRAINT __fk_fvam_assign   FOREIGN KEY     ( filterID, entity ) REFERENCES filter.assignable_entity ( filterID, entity ) DEFERRABLE,
    CONSTRAINT __fk_fvam_tlsID    FOREIGN KEY     ( tlsID ) REFERENCES ix.top_level_service ( tlsID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_prodID   FOREIGN KEY     ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_isID     FOREIGN KEY     ( isID ) REFERENCES yp.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_compID   FOREIGN KEY     ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_groupID  FOREIGN KEY     ( groupID ) REFERENCES ix.deployment_group ( groupID ) ON DELETE RESTRICT DEFERRABLE,
    CONSTRAINT __fk_fvam_orchID   FOREIGN KEY     ( orchID ) REFERENCES asset.orchestration_environment ( orchID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_rteID    FOREIGN KEY     ( rteID ) REFERENCES asset.runtime_environment ( rteID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_serverID FOREIGN KEY     ( serverID ) REFERENCES asset.server ( serverID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_endpID   FOREIGN KEY     ( endpointID ) REFERENCES ix.endpoint ( endpointID ) DEFERRABLE,
    CONSTRAINT __fk_fvam_contID   FOREIGN KEY     ( containerID ) REFERENCES asset.container ( containerID ) DEFERRABLE,
    CONSTRAINT __fk_createdBy     FOREIGN KEY     ( createdBy ) REFERENCES inventory.user ( userID ),
    CONSTRAINT __validFrom_utc    CHECK           ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
    CONSTRAINT __validUntil_utc   CHECK           ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
    CONSTRAINT __createdAt_utc    CHECK           ( EXTRACT( TIMEZONE FROM createdAt ) = '0' ),
    CONSTRAINT __fvam_cardinality CHECK           ( cardinality = 'many'::flt_card ),
    CONSTRAINT __fvam_uniq_object CHECK           (   ((entity='top_level_service')         AND (tlsID       IS NOT NULL))
                                                   OR ((entity='product')                   AND (productID   IS NOT NULL))
                                                   OR ((entity='information_system')        AND (isID        IS NOT NULL))
                                                   OR ((entity='functional_component')      AND (componentID IS NOT NULL))
                                                   OR ((entity='deployment_group')          AND (groupID     IS NOT NULL))
                                                   OR ((entity='orchestration_environment') AND (orchID      IS NOT NULL))
                                                   OR ((entity='runtime_environment')       AND (rteID       IS NOT NULL))
                                                   OR ((entity='server')                    AND (serverID    IS NOT NULL))
                                                   OR ((entity='endpoint')                  AND (endpointID  IS NOT NULL))
                                                   OR ((entity='container')                 AND (containerID IS NOT NULL))),
    CONSTRAINT __fvam_temporal    EXCLUDE         USING gist (public.uuid_to_bytea(filterValueID) WITH =,
                                                              public.uuid_to_bytea(tlsID) WITH =,
                                                              public.uuid_to_bytea(productID) WITH =,
                                                              public.uuid_to_bytea(isID) WITH =,
                                                              public.uuid_to_bytea(componentID) WITH =,
                                                              public.uuid_to_bytea(groupID) WITH =,
                                                              public.uuid_to_bytea(orchID) WITH =,
                                                              public.uuid_to_bytea(rteID) WITH =,
                                                              public.uuid_to_bytea(serverID) WITH =,
                                                              public.uuid_to_bytea(endpointID) WITH =,
                                                              public.uuid_to_bytea(containerID) WITH =,
                                                              validity WITH &&)
);
