--
--
-- VIEW SCHEMA
-- -- resolveRuntimeToServer tracks a specified runtime down across
-- -- nested/linked runtime and orchestration environments to the
-- -- next (!) server(s), which are either virtual or physical.
-- -- It does not drill further into found server(s).
CREATE FUNCTION view.resolveRuntimeToServer(rt uuid)
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
  AS (
    -- initial static anchor query
    SELECT  null::uuid,
            arep.rteID,
            null::uuid,
            arep.parentServerID,
            arep.parentRuntimeID,
            arep.parentOrchestrationID,
            0::smallint
    FROM  asset.runtime_environment_parent AS arep
    WHERE arep.rteid = rt::uuid
       OR rteID IN (
          SELECT  rteID_A
          FROM    asset.runtime_environment_linking
          WHERE   rteID_A = rt::uuid
             OR   rteID_B = rt::uuid
          UNION
          SELECT  rteID_B
          FROM    asset.runtime_environment_linking
          WHERE   rteID_A = rt::uuid
             OR   rteID_B = rt::uuid
       )
    UNION
    -- recursive iteration query
    SELECT  null::uuid,
            CASE WHEN arep.rteID  IS NOT NULL THEN arep.rteID
                 ELSE null::uuid
            END,
            CASE WHEN aoep.orchID IS NOT NULL THEN aoep.orchID
                 ELSE null::uuid
            END,
            arep.parentServerID,
            CASE WHEN arep.parentRuntimeID IS NOT NULL THEN arep.parentRuntimeID
                 WHEN aoep.parentRuntimeID IS NOT NULL THEN aoep.parentRuntimeID
                 ELSE null::uuid
            END,
            arep.parentOrchestrationID,
            t.depth+1::smallint
    FROM    t
    LEFT    JOIN
            asset.runtime_environment_parent AS arep
      ON    t.parentRuntimeID = arep.rteID
              OR arep.rteID IN (
                SELECT  rteID_A
                FROM    asset.runtime_environment_linking
                WHERE   rteID_A = t.parentruntimeid
                   OR   rteID_B = t.parentruntimeid
                UNION
                SELECT  rteID_B
                FROM    asset.runtime_environment_linking
                WHERE   rteID_A = t.parentruntimeid
                   OR   rteID_B = t.parentruntimeID
              )
     LEFT   JOIN
            asset.orchestration_environment_parent AS aoep
       ON   t.parentOrchestrationID = aoep.orchID
              OR aoep.orchID IN (
                SELECT  orchID_A
                FROM    asset.orchestration_environment_linking
                WHERE   orchID_A = t.parentOrchestrationID
                  OR    orchID_B = t.parentOrchestrationID
                UNION
                SELECT  orchID_B
                FROM    asset.orchestration_environment_linking
                WHERE   orchID_A = t.parentOrchestrationID
                   OR   orchID_B = t.parentOrchestrationID
              )
     WHERE  t.depth < 32
     )
  SELECT  ssa.serverID AS serverID,
          ssa.value    AS serverType,
          t.depth      AS depth
  FROM    asset.server_standard_attribute_values AS ssa
  JOIN    t
    ON    t.parentServerID = ssa.serverID
  NATURAL JOIN meta.standard_attribute AS ma
  WHERE   t.parentServerID IS NOT NULL
    AND   ma.attribute = 'type';
  $BODY$
  LANGUAGE sql IMMUTABLE;

-- -- resolveRuntimeToPhysical tracks a specified runtime down to the
-- -- physical server(s), across any nested virtual servers and
-- -- orchestration environments in between.
CREATE FUNCTION view.resolveRuntimeToPhysical(rt uuid)
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
  AS (
    -- initial static anchor query
    SELECT  null::uuid,
            arep.rteID,
            null::uuid,
            arep.parentServerID,
            arep.parentRuntimeID,
            arep.parentOrchestrationID,
            0::smallint
    FROM  asset.runtime_environment_parent AS arep
    WHERE arep.rteid = rt::uuid
       OR rteID IN (
          SELECT  rteID_A
          FROM    asset.runtime_environment_linking
          WHERE   rteID_A = rt::uuid
             OR   rteID_B = rt::uuid
          UNION
          SELECT  rteID_B
          FROM    asset.runtime_environment_linking
          WHERE   rteID_A = rt::uuid
             OR   rteID_B = rt::uuid
       )
    UNION
    -- recursive iteration query
    SELECT  CASE WHEN asp.serverID IS NOT NULL THEN asp.serverID
                 ELSE null::uuid
            END,
            CASE WHEN arep.rteID   IS NOT NULL THEN arep.rteID
                 ELSE null::uuid
            END,
            CASE WHEN aoep.orchID  IS NOT NULL THEN aoep.orchID
                 ELSE null::uuid
            END,
            arep.parentServerID,
            CASE WHEN arep.parentRuntimeID IS NOT NULL THEN arep.parentRuntimeID
                 WHEN  asp.parentRuntimeID IS NOT NULL THEN  asp.parentRuntimeID
                 WHEN aoep.parentRuntimeID IS NOT NULL THEN aoep.parentRuntimeID
                 ELSE null::uuid
            END,
            arep.parentOrchestrationID,
            t.depth+1::smallint
    FROM    t
    LEFT    JOIN
            asset.runtime_environment_parent AS arep
      ON    t.parentRuntimeID = arep.rteID
              OR arep.rteID IN (
                SELECT  rteID_A
                FROM    asset.runtime_environment_linking
                WHERE   rteID_A = t.parentruntimeid
                   OR   rteID_B = t.parentruntimeid
                UNION
                SELECT  rteID_B
                FROM    asset.runtime_environment_linking
                WHERE   rteID_A = t.parentruntimeid
                   OR   rteID_B = t.parentruntimeID
              )
     LEFT   JOIN
            asset.server_parent AS asp
       ON   t.parentServerID = asp.serverID
              OR asp.serverID IN (
                SELECT  serverID_A
                FROM    asset.server_linking
                WHERE   serverID_A = t.parentServerID
                   OR   serverID_B = t.parentServerID
                UNION
                SELECT  serverID_B
                FROM    asset.server_linking
                WHERE   serverID_A = t.parentServerID
                   OR   serverID_B = t.parentServerID
              )
     LEFT   JOIN
            asset.orchestration_environment_parent AS aoep
       ON   t.parentOrchestrationID = aoep.orchID
              OR aoep.orchID IN (
                SELECT  orchID_A
                FROM    asset.orchestration_environment_linking
                WHERE   orchID_A = t.parentOrchestrationID
                  OR    orchID_B = t.parentOrchestrationID
                UNION
                SELECT  orchID_B
                FROM    asset.orchestration_environment_linking
                WHERE   orchID_A = t.parentOrchestrationID
                   OR   orchID_B = t.parentOrchestrationID
              )
     WHERE  t.depth < 32
     )
  SELECT  ssa.serverID AS serverID,
          ssa.value    AS serverType,
          t.depth      AS depth
  FROM    asset.server_standard_attribute_values AS ssa
  JOIN    t
    ON    t.parentServerID = ssa.serverID
  NATURAL JOIN meta.standard_attribute AS ma
  WHERE   t.parentServerID IS NOT NULL
    AND   ma.attribute = 'type'
    AND   ssa.value = 'physical';
  $BODY$
  LANGUAGE sql IMMUTABLE;

--  EXAMPLE
--
--
--  Physical Server
--    Both resolveRuntimeToServer and resolveRuntimeToPhysical
--    return the same result since the first encountered server
--    is physical.
-- ix=> SELECT * FROM view.resolveRuntimeToServer('0af6e9a9-47b7-4d2e-82ae-7c36d05d735c');
--                serverID               | serverType | depth
-- --------------------------------------+------------+-------
--  4b64e5a9-3ac9-4ec7-a5d4-c516d0dd8077 | physical   |     1
--
-- ix=> SELECT * FROM view.resolveRuntimeToPhysical('0af6e9a9-47b7-4d2e-82ae-7c36d05d735c');
--                serverID               | serverType | depth
-- --------------------------------------+------------+-------
--  4b64e5a9-3ac9-4ec7-a5d4-c516d0dd8077 | physical   |     1
--
--  Virtual Server
--    resolveRuntimeToServer and resolveRuntimeToPhysical return
--    different results at different depths in the stack since
--    the ...ToPhysical has a drill down into virtual servers
--    that ...ToServer has not.
-- ix=> SELECT * FROM view.resolveRuntimeToServer('9dc284a6-14d8-435a-96b8-53ff7e358f0d');
--                serverID               | serverType | depth
-- --------------------------------------+------------+-------
--  8aed99bf-9ebc-4fc6-a7fb-e0f64526be01 | virtual    |     0
--
-- ix=> SELECT * FROM view.resolveRuntimeToPhysical('9dc284a6-14d8-435a-96b8-53ff7e358f0d');
--                serverID               | serverType | depth
-- --------------------------------------+------------+-------
--  4b64e5a9-3ac9-4ec7-a5d4-c516d0dd8077 | physical   |     3
--
