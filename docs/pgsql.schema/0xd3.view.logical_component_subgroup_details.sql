--
--
-- VIEW SCHEMA
CREATE  VIEW view.deployment_group_details AS
SELECT  ixdg.groupID AS groupID,
        md.dictionaryID AS groupDictionaryID,
        md.name AS groupDictionaryName,
        ixfcp.componentID AS componentID,
        ypisp.isID AS informationSystemID,
        mqa.attributeID AS attributeID,
        mqa.attribute AS attributeName,
        ixdgqv.value AS attributeValue
FROM    ix.deployment_group AS ixdg
JOIN    meta.dictionary AS md
  ON    ixdg.dictionaryID = md.dictionaryID
JOIN    ix.deployment_group_unique_attribute_values AS ixdgqv
  ON    ixdg.groupID = ixdgqv.groupID
JOIN    meta.unique_attribute AS mqa
  ON    ixdgqv.attributeID = mqa.attributeID
JOIN    ix.functional_component_parent AS ixfcp
  ON    ixdg.groupID = ixfcp.groupID
JOIN    yp.information_system_parent AS ypisp
  ON    ixfcp.componentID = ypisp.componentID
WHERE   NOW()::timestamptz(3) <@ ixdgqv.validity
  AND   NOW()::timestamptz(3) <@ ixfcp.validity
  AND   NOW()::timestamptz(3) <@ ypisp.validity
UNION
SELECT  ixdg.groupID AS groupID,
        md.dictionaryID AS groupDictionaryID,
        md.name AS groupDictionaryName,
        ixfcp.componentID AS componentID,
        ypisp.isID AS informationSystemID,
        mqa.attributeID AS attributeID,
        mqa.attribute AS attributeName,
        ixdgqv.value AS attributeValue
FROM    ix.deployment_group AS ixdg
JOIN    meta.dictionary AS md
  ON    ixdg.dictionaryID = md.dictionaryID
JOIN    ix.deployment_group_unique_attribute_values AS ixdgqv
  ON    ixdg.groupID = ixdgqv.groupID
JOIN    meta.unique_attribute AS mqa
  ON    ixdgqv.attributeID = mqa.attributeID
JOIN    ix.functional_component_parent AS ixfcp
  ON    ixdg.groupID = ixfcp.groupID
JOIN    yp.information_system_parent AS ypisp
  ON    ixfcp.componentID = ypisp.componentID
WHERE   NOW()::timestamptz(3) <@ ixdgqv.validity
  AND   NOW()::timestamptz(3) <@ ixfcp.validity
  AND   NOW()::timestamptz(3) <@ ypisp.validity;

CREATE  FUNCTION view.deployment_group_details_at(at timestamptz)
  RETURNS TABLE ( groupID             uuid,
                  groupDictionaryID   uuid,
                  groupDictionaryName text,
                  componentID         uuid,
                  informationSystemID uuid,
                  attributeID         uuid,
                  attributeName       text,
                  attributeValue      text)
  AS
  $BODY$
  SELECT  ixdg.groupID AS groupID,
          md.dictionaryID AS groupDictionaryID,
          md.name AS groupDictionaryName,
          ixfcp.componentID AS componentID,
          ypisp.isID AS informationSystemID,
          mqa.attributeID AS attributeID,
          mqa.attribute AS attributeName,
          ixdgqv.value AS attributeValue
  FROM    ix.deployment_group AS ixdg
  JOIN    meta.dictionary AS md
    ON    ixdg.dictionaryID = md.dictionaryID
  JOIN    ix.deployment_group_unique_attribute_values AS ixdgqv
    ON    ixdg.groupID = ixdgqv.groupID
  JOIN    meta.unique_attribute AS mqa
    ON    ixdgqv.attributeID = mqa.attributeID
  JOIN    ix.functional_component_parent AS ixfcp
    ON    ixdg.groupID = ixfcp.groupID
  JOIN    yp.information_system_parent AS ypisp
    ON    ixfcp.componentID = ypisp.componentID
  WHERE   at::timestamptz(3) <@ ixdgqv.validity
    AND   at::timestamptz(3) <@ ixfcp.validity
    AND   at::timestamptz(3) <@ ypisp.validity
  UNION
  SELECT  ixdg.groupID AS groupID,
          md.dictionaryID AS groupDictionaryID,
          md.name AS groupDictionaryName,
          ixfcp.componentID AS componentID,
          ypisp.isID AS informationSystemID,
          mqa.attributeID AS attributeID,
          mqa.attribute AS attributeName,
          ixdgqv.value AS attributeValue
  FROM    ix.deployment_group AS ixdg
  JOIN    meta.dictionary AS md
    ON    ixdg.dictionaryID = md.dictionaryID
  JOIN    ix.deployment_group_unique_attribute_values AS ixdgqv
    ON    ixdg.groupID = ixdgqv.groupID
  JOIN    meta.unique_attribute AS mqa
    ON    ixdgqv.attributeID = mqa.attributeID
  JOIN    ix.functional_component_parent AS ixfcp
    ON    ixdg.groupID = ixfcp.groupID
  JOIN    yp.information_system_parent AS ypisp
    ON    ixfcp.componentID = ypisp.componentID
  WHERE   at::timestamptz(3) <@ ixdgqv.validity
    AND   at::timestamptz(3) <@ ixfcp.validity
    AND   at::timestamptz(3) <@ ypisp.validity
  $BODY$
  LANGUAGE sql IMMUTABLE;
