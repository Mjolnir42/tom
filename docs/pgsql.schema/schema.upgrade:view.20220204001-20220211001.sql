BEGIN;

  CREATE OR REPLACE FUNCTION view.resolveRuntimeToServer(rt uuid)
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
            WHERE   rteID_B = rt::uuid
            UNION
            SELECT  rteID_B
            FROM    asset.runtime_environment_linking
            WHERE   rteID_A = rt::uuid
         )
      UNION
      -- recursive iteration query
      SELECT  CASE WHEN t.parentServerID IS NOT NULL
                   THEN t.parentServerID
                   ELSE null::uuid
              END,
              CASE WHEN arep.rteID  IS NOT NULL THEN arep.rteID
                   ELSE null::uuid
              END,
              CASE WHEN aoem.orchID IS NOT NULL THEN aoem.orchID
                   ELSE null::uuid
              END,
              arep.parentServerID,
              CASE WHEN arep.parentRuntimeID IS NOT NULL THEN arep.parentRuntimeID
                   WHEN aoem.parentRuntimeID IS NOT NULL THEN aoem.parentRuntimeID
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
                  WHERE   rteID_B = t.parentruntimeid
                  UNION
                  SELECT  rteID_B
                  FROM    asset.runtime_environment_linking
                  WHERE   rteID_A = t.parentruntimeid
                )
       LEFT   JOIN
              asset.orchestration_environment_mapping AS aoem
         ON   t.parentOrchestrationID = aoem.orchID
                OR aoem.orchID IN (
                  SELECT  orchID_A
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_B = t.parentOrchestrationID
                  UNION
                  SELECT  orchID_B
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_A = t.parentOrchestrationID
                )
       WHERE  t.depth < 32
         AND  (   t.parentRuntimeID       IS NOT NULL
               OR t.parentServerID        IS NOT NULL
               OR t.parentOrchestrationID IS NOT NULL)
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
      AND   ma.attribute = 'type';
    $BODY$
    LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveRuntimeToPhysical(rt uuid)
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
            WHERE   rteID_B = rt::uuid
            UNION
            SELECT  rteID_B
            FROM    asset.runtime_environment_linking
            WHERE   rteID_A = rt::uuid
         )
      UNION
      -- recursive iteration query
      SELECT  CASE WHEN asp.serverID IS NOT NULL THEN asp.serverID
                   WHEN t.parentServerID IS NOT NULL AND asp.serverID IS NULL THEN t.parentServerID
                   ELSE null::uuid
              END,
              CASE WHEN arep.rteID   IS NOT NULL THEN arep.rteID
                   ELSE null::uuid
              END,
              CASE WHEN aoem.orchID  IS NOT NULL THEN aoem.orchID
                   ELSE null::uuid
              END,
              arep.parentServerID,
              CASE WHEN arep.parentRuntimeID IS NOT NULL THEN arep.parentRuntimeID
                   WHEN  asp.parentRuntimeID IS NOT NULL THEN  asp.parentRuntimeID
                   WHEN aoem.parentRuntimeID IS NOT NULL THEN aoem.parentRuntimeID
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
                  WHERE   rteID_B = t.parentRuntimeID
                  UNION
                  SELECT  rteID_B
                  FROM    asset.runtime_environment_linking
                  WHERE   rteID_A = t.parentRuntimeID
                )
       LEFT   JOIN
              asset.server_parent AS asp
         ON   t.parentServerID = asp.serverID
                OR asp.serverID IN (
                  SELECT  serverID_A
                  FROM    asset.server_linking
                  WHERE   serverID_B = t.parentServerID
                  UNION
                  SELECT  serverID_B
                  FROM    asset.server_linking
                  WHERE   serverID_A = t.parentServerID
                )
       LEFT   JOIN
              asset.orchestration_environment_mapping AS aoem
         ON   t.parentOrchestrationID = aoem.orchID
                OR aoem.orchID IN (
                  SELECT  orchID_A
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_B = t.parentOrchestrationID
                  UNION
                  SELECT  orchID_B
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_A = t.parentOrchestrationID
                )
       WHERE  t.depth < 32
         AND  (   t.parentRuntimeID       IS NOT NULL
               OR t.parentServerID        IS NOT NULL
               OR t.parentOrchestrationID IS NOT NULL)
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
      AND   ssa.value = 'physical';
    $BODY$
    LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveContainerToServer(cnID uuid)
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
      WHERE asset.container_parent.containerID = cnID::uuid
         OR asset.container_parent.containerID IN (
            SELECT  containerID_A
            FROM    asset.container_linking
            WHERE  asset.container_linking.containerID_B = cnID::uuid
            UNION
            SELECT  containerID_B
            FROM    asset.container_linking
            WHERE  asset.container_linking.containerID_A = cnID::uuid
         )
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
      AND   ma.attribute = 'type';
  $BODY$
  LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveContainerToPhysical(cnID uuid)
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
      WHERE asset.container_parent.containerID = cnID::uuid
         OR asset.container_parent.containerID IN (
            SELECT containerID_A FROM asset.container_linking
            WHERE  asset.container_linking.containerID_B = cnID::uuid
            UNION
            SELECT containerID_B FROM asset.container_linking
            WHERE  asset.container_linking.containerID_A = cnID::uuid)
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
        AND   ssa.value = 'physical';
  $BODY$
  LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveOrchestrationToServer(oreID uuid)
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
              oreID::uuid,
              null::uuid,
              asset.orchestration_environment_mapping.parentRuntimeID,
              null::uuid,
              0::smallint
      FROM  asset.orchestration_environment_mapping
      WHERE asset.orchestration_environment_mapping.orchID = oreID::uuid
         OR asset.orchestration_environment_mapping.orchID IN (
            SELECT  orchID_A
            FROM    asset.orchestration_environment_linking
            WHERE  asset.orchestration_environment_linking.orchID_B = oreID::uuid
            UNION
            SELECT  orchID_B
            FROM    asset.orchestration_environment_linking
            WHERE  asset.orchestration_environment_linking.orchID_A = oreID::uuid
         )
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
      LEFT    JOIN
              asset.runtime_environment_parent
        ON    t.parentRuntimeID = asset.runtime_environment_parent.rteID
                OR asset.runtime_environment_parent.rteID IN (
                  SELECT  rteID_A
                  FROM    asset.runtime_environment_linking
                  WHERE   rteID_B = t.parentruntimeid
                  UNION
                  SELECT  rteID_B
                  FROM    asset.runtime_environment_linking
                  WHERE   rteID_A = t.parentruntimeid
                )
       LEFT   JOIN
              asset.orchestration_environment_mapping
         ON   t.parentOrchestrationID = asset.orchestration_environment_mapping.orchID
                OR asset.orchestration_environment_mapping.orchID IN (
                  SELECT  orchID_A
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_B = t.parentOrchestrationID
                  UNION
                  SELECT  orchID_B
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_A = t.parentOrchestrationID
                )
       WHERE  t.depth < 32
    AND  ( t.parentRuntimeID IS NOT NULL OR
           t.parentServerID IS NOT NULL OR
           t.parentOrchestrationID IS NOT NULL )
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
      AND   ma.attribute = 'type';
    $BODY$
    LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveOrchestrationToPhysical(oreID uuid)
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
                    depth) AS ( -- initial static anchor query
  SELECT null::uuid,
         null::uuid,
         oreID::uuid,
         null::uuid,
         asset.orchestration_environment_mapping.parentRuntimeID,
         null::uuid,
         0::smallint
  FROM   asset.orchestration_environment_mapping
  WHERE  asset.orchestration_environment_mapping.orchID = oreID::uuid
     OR  asset.orchestration_environment_mapping.orchID IN (
         SELECT orchID_A
         FROM   asset.orchestration_environment_linking
         WHERE  asset.orchestration_environment_linking.orchID_B = oreID::uuid
         UNION
         SELECT orchID_B
         FROM   asset.orchestration_environment_linking
         WHERE  asset.orchestration_environment_linking.orchID_A = oreID::uuid
     )
  UNION -- recursive iteration query
  SELECT CASE WHEN asset.server_parent.serverID IS NOT NULL THEN asset.server_parent.serverID
              WHEN t.parentServerID IS NOT NULL AND asset.server_parent.serverID IS NULL THEN t.parentServerID
                   ELSE null::uuid
         END,
         CASE WHEN asset.runtime_environment_parent.rteID IS NOT NULL THEN asset.runtime_environment_parent.rteID
                   ELSE null::uuid
         END,
         CASE WHEN asset.orchestration_environment_mapping.orchID IS NOT NULL THEN  asset.orchestration_environment_mapping.orchID
                   ELSE null::uuid
         END,
         asset.runtime_environment_parent.parentServerID,
         CASE WHEN asset.runtime_environment_parent.parentRuntimeID IS NOT NULL THEN asset.runtime_environment_parent.parentRuntimeID
              WHEN asset.server_parent.parentRuntimeID IS NOT NULL THEN asset.server_parent.parentRuntimeID
              WHEN asset.orchestration_environment_mapping.parentRuntimeID IS NOT NULL THEN asset.orchestration_environment_mapping.parentRuntimeID
              ELSE null::uuid
         END,
         asset.runtime_environment_parent.parentOrchestrationID,
         t.depth+1::smallint
  FROM   t
  LEFT   JOIN
         asset.runtime_environment_parent
      ON t.parentRuntimeID = asset.runtime_environment_parent.rteID
         OR asset.runtime_environment_parent.rteID IN (
            SELECT  rteID_A FROM    asset.runtime_environment_linking
            WHERE rteID_B = t.parentruntimeid
            UNION
            SELECT  rteID_B FROM    asset.runtime_environment_linking
            WHERE   rteID_A = t.parentruntimeid )
  LEFT   JOIN
         asset.server_parent
      ON t.parentServerID = asset.server_parent.serverID
         OR asset.server_parent.serverID IN (
            SELECT serverID_A FROM asset.server_linking
            WHERE  serverID_B = t.parentServerID
            UNION
            SELECT serverID_B FROM asset.server_linking
            WHERE  serverID_A = t.parentServerID )
  LEFT   JOIN
         asset.orchestration_environment_mapping
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
  )

  SELECT ssa.serverID AS serverID,
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
      AND   ssa.value = 'physical';

    $BODY$
    LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveServerToServer(srvID uuid)
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
      SELECT  srvID::uuid,
              null::uuid,
              null::uuid,
              null::uuid,
              null::uuid,
              null::uuid,
              0::smallint
      FROM  asset.server
      LEFT  JOIN asset.server_parent
         ON asset.server.serverID = asset.server_parent.serverID
      WHERE asset.server.serverID = srvID::uuid
         OR asset.server_parent.serverID IN (
            SELECT  serverID_A
            FROM    asset.server_linking
            WHERE  asset.server_linking.serverID_B = srvID::uuid
            UNION
            SELECT  serverID_B
            FROM    asset.server_linking
            WHERE  asset.server_linking.serverID_A = srvID::uuid
         )
      UNION
      -- recursive iteration query
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
      LEFT    JOIN
              asset.runtime_environment_parent
        ON    t.parentRuntimeID = asset.runtime_environment_parent.rteID
                OR asset.runtime_environment_parent.rteID IN (
                  SELECT  rteID_A
                  FROM    asset.runtime_environment_linking
                  WHERE   rteID_B = t.parentruntimeid
                  UNION
                  SELECT  rteID_B
                  FROM    asset.runtime_environment_linking
                  WHERE   rteID_A = t.parentruntimeid
                )
       LEFT   JOIN
              asset.orchestration_environment_mapping
         ON   t.parentOrchestrationID = asset.orchestration_environment_mapping.orchID
                OR asset.orchestration_environment_mapping.orchID IN (
                  SELECT  orchID_A
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_B = t.parentOrchestrationID
                  UNION
                  SELECT  orchID_B
                  FROM    asset.orchestration_environment_linking
                  WHERE   orchID_A = t.parentOrchestrationID
                )
       WHERE  t.depth < 32
    AND  ( t.parentRuntimeID IS NOT NULL OR
           t.parentServerID IS NOT NULL OR
           t.parentOrchestrationID IS NOT NULL )
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
      AND   ma.attribute = 'type';
    $BODY$
    LANGUAGE sql IMMUTABLE;

  CREATE OR REPLACE FUNCTION view.resolveServerToPhysical(srvID uuid)
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
                    depth) AS (
                    

  -- initial static anchor query
  SELECT srvID::uuid,
         null::uuid,
         null::uuid,
         null::uuid,
         asset.server_parent.parentRuntimeID,
         null::uuid,
         0::smallint
  FROM   asset.server
  LEFT   JOIN asset.server_parent
     ON  asset.server.serverID = asset.server_parent.serverID
  WHERE  asset.server.serverID = srvID::uuid
     OR  asset.server_parent.serverID IN (
         SELECT serverID_A
         FROM   asset.server_linking
         WHERE  asset.server_linking.serverID_B = srvID::uuid
         UNION
         SELECT serverID_B
         FROM   asset.server_linking
         WHERE  asset.server_linking.serverID_A = srvID::uuid
     )
  --
  UNION
  -- recursive iteration query
  SELECT CASE WHEN asset.server_parent.serverID IS NOT NULL THEN asset.server_parent.serverID
              WHEN t.parentServerID IS NOT NULL AND asset.server_parent.serverID IS NULL THEN t.parentServerID
                   ELSE null::uuid
         END,
         CASE WHEN asset.runtime_environment_parent.rteID IS NOT NULL THEN asset.runtime_environment_parent.rteID
                   ELSE null::uuid
         END,
         CASE WHEN asset.orchestration_environment_mapping.orchID IS NOT NULL THEN  asset.orchestration_environment_mapping.orchID
                   ELSE null::uuid
         END,
         asset.runtime_environment_parent.parentServerID,
         CASE WHEN asset.runtime_environment_parent.parentRuntimeID IS NOT NULL THEN asset.runtime_environment_parent.parentRuntimeID
              WHEN asset.server_parent.parentRuntimeID IS NOT NULL THEN asset.server_parent.parentRuntimeID
              WHEN asset.orchestration_environment_mapping.parentRuntimeID IS NOT NULL THEN asset.orchestration_environment_mapping.parentRuntimeID
              ELSE null::uuid
         END,
         asset.runtime_environment_parent.parentOrchestrationID,
         t.depth+1::smallint
  FROM   t
  LEFT   JOIN
         asset.runtime_environment_parent
      ON t.parentRuntimeID = asset.runtime_environment_parent.rteID
         OR asset.runtime_environment_parent.rteID IN (
            SELECT  rteID_A FROM    asset.runtime_environment_linking
            WHERE rteID_B = t.parentruntimeid
            UNION
            SELECT  rteID_B FROM    asset.runtime_environment_linking
            WHERE   rteID_A = t.parentruntimeid )
  LEFT   JOIN
         asset.server_parent
      ON t.parentServerID = asset.server_parent.serverID
         OR asset.server_parent.serverID IN (
            SELECT serverID_A FROM asset.server_linking
            WHERE  serverID_B = t.parentServerID
            UNION
            SELECT serverID_B FROM asset.server_linking
            WHERE  serverID_A = t.parentServerID )
  LEFT   JOIN
         asset.orchestration_environment_mapping
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
  )

  SELECT ssa.serverID AS serverID,
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
      AND   ssa.value = 'physical';

    $BODY$
    LANGUAGE sql IMMUTABLE;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'view', 20220211001, 'update view.resolve..To.. functions');
COMMIT;
