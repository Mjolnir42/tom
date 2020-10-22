BEGIN;
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



  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'view', 20191104002, 'add resolveRuntimeTo.. functions');
COMMIT;
