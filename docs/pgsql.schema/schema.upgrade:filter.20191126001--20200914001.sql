-- ALTER TYPE can not run inside a transaction
ALTER TYPE flt_ntt  ADD VALUE 'endpoint' AFTER 'server';

BEGIN;
  ALTER TABLE filter.value_assignment__one ADD COLUMN endpointID uuid NULL;
  ALTER TABLE filter.value_assignment__one ADD CONSTRAINT __fk_fvao_endpID FOREIGN KEY ( endpointID ) REFERENCES ix.endpoint ( endpointID ) DEFERRABLE;
  ALTER TABLE filter.value_assignment__one DROP CONSTRAINT __fvao_uniq_object;
  ALTER TABLE filter.value_assignment__one DROP CONSTRAINT __fvao_temporal;
  ALTER TABLE filter.value_assignment__one ADD CONSTRAINT __fvao_uniq_object CHECK
                                              (   ((entity='top_level_service')         AND (tlsID       IS NOT NULL))
                                               OR ((entity='product')                   AND (productID   IS NOT NULL))
                                               OR ((entity='information_system')        AND (isID        IS NOT NULL))
                                               OR ((entity='functional_component')      AND (componentID IS NOT NULL))
                                               OR ((entity='deployment_group')          AND (groupID     IS NOT NULL))
                                               OR ((entity='orchestration_environment') AND (orchID      IS NOT NULL))
                                               OR ((entity='runtime_environment')       AND (rteID       IS NOT NULL))
                                               OR ((entity='server')                    AND (serverID    IS NOT NULL))
                                               OR ((entity='endpoint')                  AND (endpointID  IS NOT NULL)));
  ALTER TABLE filter.value_assignment__one ADD CONSTRAINT __fvao_temporal    EXCLUDE     USING gist
                                                         (public.uuid_to_bytea(filterID) WITH =,
                                                          public.uuid_to_bytea(tlsID) WITH =,
                                                          public.uuid_to_bytea(productID) WITH =,
                                                          public.uuid_to_bytea(isID) WITH =,
                                                          public.uuid_to_bytea(componentID) WITH =,
                                                          public.uuid_to_bytea(groupID) WITH =,
                                                          public.uuid_to_bytea(orchID) WITH =,
                                                          public.uuid_to_bytea(rteID) WITH =,
                                                          public.uuid_to_bytea(serverID) WITH =,
                                                          public.uuid_to_bytea(endpointID) WITH =,
                                                          validity WITH &&);

  ALTER TABLE filter.value_assignment__many ADD COLUMN endpointID uuid NULL;
  ALTER TABLE filter.value_assignment__many ADD CONSTRAINT __fk_fvam_endpID FOREIGN KEY ( endpointID ) REFERENCES ix.endpoint ( endpointID ) DEFERRABLE;
  ALTER TABLE filter.value_assignment__many DROP CONSTRAINT __fvam_uniq_object;
  ALTER TABLE filter.value_assignment__many DROP CONSTRAINT __fvam_temporal;
  ALTER TABLE filter.value_assignment__many ADD CONSTRAINT __fvam_uniq_object CHECK
                                              (   ((entity='top_level_service')         AND (tlsID       IS NOT NULL))
                                               OR ((entity='product')                   AND (productID   IS NOT NULL))
                                               OR ((entity='information_system')        AND (isID        IS NOT NULL))
                                               OR ((entity='functional_component')      AND (componentID IS NOT NULL))
                                               OR ((entity='deployment_group')          AND (groupID     IS NOT NULL))
                                               OR ((entity='orchestration_environment') AND (orchID      IS NOT NULL))
                                               OR ((entity='runtime_environment')       AND (rteID       IS NOT NULL))
                                               OR ((entity='server')                    AND (serverID    IS NOT NULL))
                                               OR ((entity='endpoint')                  AND (endpointID  IS NOT NULL)));
  ALTER TABLE filter.value_assignment__many ADD CONSTRAINT __fvam_temporal EXCLUDE USING gist
                                                         (public.uuid_to_bytea(filterID) WITH =,
                                                          public.uuid_to_bytea(tlsID) WITH =,
                                                          public.uuid_to_bytea(productID) WITH =,
                                                          public.uuid_to_bytea(isID) WITH =,
                                                          public.uuid_to_bytea(componentID) WITH =,
                                                          public.uuid_to_bytea(groupID) WITH =,
                                                          public.uuid_to_bytea(orchID) WITH =,
                                                          public.uuid_to_bytea(rteID) WITH =,
                                                          public.uuid_to_bytea(serverID) WITH =,
                                                          public.uuid_to_bytea(endpointID) WITH =,
                                                          validity WITH &&);

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'filter', 20200914001, 'add endpoint as filter-able entity');
COMMIT;
