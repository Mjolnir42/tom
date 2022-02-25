--
--
-- VIEW SCHEMA
-- -- resolveContainerToServerAt tracks a specified container down across
-- -- nested/linked runtime and orchestration environments to the
-- -- next (!) server(s), which are either virtual or physical.
-- -- It does not drill further into found server(s).
CREATE OR REPLACE FUNCTION view.resolveContainerToServerAt(cnID uuid, at timestamptz)
  RETURNS TABLE ( serverID   uuid,
                  serverType text,
                  depth      smallint)
AS
$BODY$
  WITH RECURSIVE t(
    serverID,
    runtimeID,
    orchestrationID,
    parentServerID,
    parentRuntimeID,
    parentOrchestrationID,
    depth)
  AS ( -- initial static anchor query
    SELECT  null::uuid,
            null::uuid,
            null::uuid,
            null::uuid,
            asset.container_parent.parentRuntimeID,
            null::uuid,
            0::smallint
    FROM  asset.container_parent
    WHERE ( asset.container_parent.containerID = cnID::uuid
       OR asset.container_parent.containerID IN (
          SELECT  containerID_A
          FROM    asset.container_linking
          WHERE  asset.container_linking.containerID_B = cnID::uuid
          UNION
          SELECT  containerID_B
          FROM    asset.container_linking
          WHERE  asset.container_linking.containerID_A = cnID::uuid
       ) )
      AND at <@ asset.container_parent.validity
    UNION -- recursive iteration query
    SELECT  CASE WHEN t.parentServerID IS NOT NULL
                      THEN t.parentServerID
                 ELSE null::uuid
            END,
            CASE WHEN asset.runtime_environment_parent.rteID IS NOT NULL
                      THEN asset.runtime_environment_parent.rteID
                 ELSE null::uuid
            END,
            CASE WHEN asset.orchestration_environment_mapping.orchID IS NOT NULL
                      THEN asset.orchestration_environment_mapping.orchID
                 ELSE null::uuid
            END,
            asset.runtime_environment_parent.parentServerID,
            CASE WHEN asset.runtime_environment_parent.parentRuntimeID IS NOT NULL
                      THEN asset.runtime_environment_parent.parentRuntimeID
                 WHEN asset.orchestration_environment_mapping.parentRuntimeID IS NOT NULL
                      THEN asset.orchestration_environment_mapping.parentRuntimeID
                 ELSE null::uuid
            END,
            asset.runtime_environment_parent.parentOrchestrationID,
            t.depth+1::smallint
    FROM    t
    LEFT    JOIN asset.runtime_environment_parent
        ON  t.parentRuntimeID = asset.runtime_environment_parent.rteID
            OR asset.runtime_environment_parent.rteID IN (
              SELECT  rteID_A FROM asset.runtime_environment_linking
              WHERE   rteID_B = t.parentruntimeid
              UNION
              SELECT  rteID_B FROM asset.runtime_environment_linking
              WHERE   rteID_A = t.parentruntimeid)
     LEFT   JOIN asset.orchestration_environment_mapping
         ON t.parentOrchestrationID = asset.orchestration_environment_mapping.orchID
            OR asset.orchestration_environment_mapping.orchID IN (
              SELECT  orchID_A FROM asset.orchestration_environment_linking
              WHERE   orchID_B = t.parentOrchestrationID
              UNION
              SELECT  orchID_B FROM asset.orchestration_environment_linking
              WHERE   orchID_A = t.parentOrchestrationID)
     WHERE  t.depth < 32
       AND  (   t.parentRuntimeID       IS NOT NULL
             OR t.parentServerID        IS NOT NULL
             OR t.parentOrchestrationID IS NOT NULL )
       AND  ( at <@ asset.runtime_environment_parent.validity OR asset.runtime_environment_parent.validity IS NULL )
       AND  ( at <@ asset.orchestration_environment_mapping.validity OR asset.orchestration_environment_mapping.validity IS NULL )
  )
  SELECT  ssa.serverID AS serverID,
          ssa.value    AS serverType,
          t.depth      AS depth
  FROM    asset.server_standard_attribute_values AS ssa
  JOIN    t
    ON    t.serverID = ssa.serverID
  JOIN    meta.standard_attribute AS ma
    ON    ssa.dictionaryID = ma.dictionaryID
   AND    ssa.attributeID = ma.attributeID
  WHERE   t.serverID IS NOT NULL
    AND   t.parentServerID IS NULL
    AND   t.parentRuntimeID IS NULL
    AND   t.parentOrchestrationID IS NULL
    AND   ma.attribute = 'type'
    AND   at <@ ssa.validity;
$BODY$
LANGUAGE sql IMMUTABLE;

-- -- resolveContainerToPhysicalAt tracks a specified container down to the
-- -- physical server(s), across any nested virtual servers and
-- -- orchestration environments in between.
CREATE OR REPLACE FUNCTION view.resolveContainerToPhysicalAt(cnID uuid, at timestamptz)
  RETURNS TABLE ( serverID uuid,
                  serverType text,
                  depth      smallint)
AS
$BODY$
  WITH RECURSIVE t (
    serverID,
    runtimeID,
    orchestrationID,
    parentServerID,
    parentRuntimeID,
    parentOrchestrationID,
    depth)
  AS ( -- initial static anchor query
    SELECT null::uuid,
           null::uuid,
           null::uuid,
           null::uuid,
           asset.container_parent.parentRuntimeID,
           null::uuid,
           0::smallint
    FROM  asset.container_parent
    WHERE ( asset.container_parent.containerID = cnID::uuid
       OR asset.container_parent.containerID IN (
          SELECT containerID_A FROM asset.container_linking
          WHERE  asset.container_linking.containerID_B = cnID::uuid
          UNION
          SELECT containerID_B FROM asset.container_linking
          WHERE  asset.container_linking.containerID_A = cnID::uuid) )
       AND at <@ asset.container_parent.validity
    UNION -- recursive iteration query
    SELECT CASE WHEN asset.server_parent.serverID IS NOT NULL
                     THEN asset.server_parent.serverID
                WHEN t.parentServerID IS NOT NULL AND asset.server_parent.serverID IS NULL
                     THEN t.parentServerID
                ELSE null::uuid
           END,
           CASE WHEN asset.runtime_environment_parent.rteID IS NOT NULL
                     THEN asset.runtime_environment_parent.rteID
                ELSE null::uuid
           END,
           CASE WHEN asset.orchestration_environment_mapping.orchID IS NOT NULL
                     THEN asset.orchestration_environment_mapping.orchID
                ELSE null::uuid
           END,
           asset.runtime_environment_parent.parentServerID,
           CASE WHEN asset.runtime_environment_parent.parentRuntimeID IS NOT NULL
                     THEN asset.runtime_environment_parent.parentRuntimeID
                WHEN asset.server_parent.parentRuntimeID IS NOT NULL
                     THEN asset.server_parent.parentRuntimeID
                WHEN asset.orchestration_environment_mapping.parentRuntimeID IS NOT NULL
                     THEN asset.orchestration_environment_mapping.parentRuntimeID
                ELSE null::uuid
           END,
           asset.runtime_environment_parent.parentOrchestrationID,
           t.depth+1::smallint
    FROM   t
    LEFT   JOIN asset.runtime_environment_parent
        ON t.parentRuntimeID = asset.runtime_environment_parent.rteID
           OR asset.runtime_environment_parent.rteID IN (
              SELECT  rteID_A FROM asset.runtime_environment_linking
              WHERE rteID_B = t.parentruntimeid
              UNION
              SELECT  rteID_B FROM asset.runtime_environment_linking
              WHERE   rteID_A = t.parentruntimeid )
    LEFT   JOIN asset.server_parent
        ON t.parentServerID = asset.server_parent.serverID
           OR asset.server_parent.serverID IN (
              SELECT serverID_A FROM asset.server_linking
              WHERE  serverID_B = t.parentServerID
              UNION
              SELECT serverID_B FROM asset.server_linking
              WHERE  serverID_A = t.parentServerID )
    LEFT   JOIN asset.orchestration_environment_mapping
        ON t.parentOrchestrationID = asset.orchestration_environment_mapping.orchID
           OR asset.orchestration_environment_mapping.orchID IN (
              SELECT orchID_A FROM asset.orchestration_environment_linking
              WHERE  orchID_B = t.parentOrchestrationID
              UNION
              SELECT orchID_B FROM asset.orchestration_environment_linking
              WHERE  orchID_A = t.parentOrchestrationID )
    WHERE  t.depth < 32
      AND  ( t.parentRuntimeID IS NOT NULL OR
             t.parentServerID IS NOT NULL OR
             t.parentOrchestrationID IS NOT NULL )
      AND  ( at <@ asset.runtime_environment_parent.validity OR asset.runtime_environment_parent.validity IS NULL )
      AND  ( at <@ asset.server_parent.validity OR asset.server_parent.validity IS NULL )
      AND  ( at <@ asset.orchestration_environment_mapping.validity OR asset.orchestration_environment_mapping.validity IS NULL )
  )
  SELECT    ssa.serverID AS serverID,
            ssa.value    AS serverType,
            t.depth      AS depth
  FROM      asset.server_standard_attribute_values AS ssa
    JOIN    t
      ON    t.serverID = ssa.serverID
    JOIN    meta.standard_attribute AS ma
      ON    ssa.dictionaryID = ma.dictionaryID
     AND    ssa.attributeID  = ma.attributeID
    WHERE   t.serverID          IS NOT NULL
      AND   t.parentServerID        IS NULL
      AND   t.parentRuntimeID       IS NULL
      AND   t.parentOrchestrationID IS NULL
      AND   ma.attribute = 'type'
      AND   ssa.value = 'physical'
      AND   at <@ ssa.validity;
$BODY$
LANGUAGE sql IMMUTABLE;
