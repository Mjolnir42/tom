BEGIN;
  DROP VIEW view.logical_component_subgroup_details;
  CREATE VIEW view.deployment_group_details
  AS          SELECT ixdg.groupID AS groupID, md.dictionaryID AS groupDictionaryID, md.name AS groupDictionaryName,
              ixmfc.componentID AS componentID, ixmis.isID AS informationSystemID, mqa.attributeID AS attributeID, mqa.attribute AS attributeName,
              ixdgqv.value AS attributeValue FROM ix.deployment_group AS ixdg JOIN meta.dictionary AS md ON ixdg.dictionaryID = md.dictionaryID
              JOIN ix.deployment_group_attr_uniq_values AS ixdgqv ON ixdg.groupID = ixdgqv.groupID JOIN meta.unique_attribute AS mqa ON
              ixdgqv.attributeID = mqa.attributeID JOIN ix.mapping_functional_component AS ixmfc ON ixdg.groupID = ixmfc.groupID JOIN
              ix.mapping_information_system AS ixmis ON ixmfc.componentID = ixmis.componentID WHERE NOW()::timestamptz(3) <@ ixdgqv.validity
              AND NOW()::timestamptz(3) <@ ixmfc.validity AND NOW()::timestamptz(3) <@ ixmis.validity
  UNION       SELECT ixdg.groupID AS groupID, md.dictionaryID AS groupDictionaryID, md.name AS groupDictionaryName,
              ixmfc.componentID AS componentID, ixmis.isID AS informationSystemID, mqa.attributeID AS attributeID, mqa.attribute AS attributeName,
              ixdgqv.value AS attributeValue FROM ix.deployment_group AS ixdg JOIN meta.dictionary AS md ON ixdg.dictionaryID = md.dictionaryID
              JOIN ix.deployment_group_attr_uniq_values AS ixdgqv ON ixdg.groupID = ixdgqv.groupID JOIN meta.unique_attribute AS mqa
              ON ixdgqv.attributeID = mqa.attributeID JOIN ix.mapping_functional_component AS ixmfc ON ixdg.groupID = ixmfc.groupID JOIN
              ix.mapping_information_system AS ixmis ON ixmfc.componentID = ixmis.componentID WHERE NOW()::timestamptz(3) <@ ixdgqv.validity
              AND NOW()::timestamptz(3) <@ ixmfc.validity AND NOW()::timestamptz(3) <@ ixmis.validity;

  DROP FUNCTION IF EXISTS view.logical_component_subgroup_details_at;
  CREATE  FUNCTION view.deployment_group_details_at(at timestamptz)
  RETURNS TABLE ( groupID             uuid, groupDictionaryID   uuid, groupDictionaryName text, componentID         uuid,
                  informationSystemID uuid, attributeID         uuid, attributeName       text, attributeValue      text)
  AS
  $BODY$
  SELECT  ixdg.groupID AS groupID, md.dictionaryID AS groupDictionaryID, md.name AS groupDictionaryName, ixmfc.componentID AS componentID,
          ixmis.isID AS informationSystemID, mqa.attributeID AS attributeID, mqa.attribute AS attributeName, ixdgqv.value AS attributeValue
          FROM ix.deployment_group AS ixdg JOIN meta.dictionary AS md ON ixdg.dictionaryID = md.dictionaryID JOIN ix.deployment_group_attr_uniq_values
          AS ixdgqv ON ixdg.groupID = ixdgqv.groupID JOIN meta.unique_attribute AS mqa ON ixdgqv.attributeID = mqa.attributeID JOIN
          ix.mapping_functional_component AS ixmfc ON ixdg.groupID = ixmfc.groupID JOIN ix.mapping_information_system AS ixmis ON
          ixmfc.componentID = ixmis.componentID WHERE at::timestamptz(3) <@ ixdgqv.validity AND at::timestamptz(3) <@ ixmfc.validity
          AND at::timestamptz(3) <@ ixmis.validity
  UNION   SELECT ixdg.groupID AS groupID, md.dictionaryID AS groupDictionaryID, md.name AS groupDictionaryName, ixmfc.componentID AS componentID,
          ixmis.isID AS informationSystemID, mqa.attributeID AS attributeID, mqa.attribute AS attributeName, ixdgqv.value AS attributeValue FROM ix.deployment_group AS ixdg
          JOIN meta.dictionary AS md ON ixdg.dictionaryID = md.dictionaryID JOIN ix.deployment_group_attr_uniq_values AS ixdgqv ON ixdg.groupID = ixdgqv.groupID
          JOIN meta.unique_attribute AS mqa ON ixdgqv.attributeID = mqa.attributeID JOIN ix.mapping_functional_component AS ixmfc ON ixdg.groupID = ixmfc.groupID
          JOIN ix.mapping_information_system AS ixmis ON ixmfc.componentID = ixmis.componentID WHERE at::timestamptz(3) <@ ixdgqv.validity AND at::timestamptz(3) <@ ixmfc.validity
          AND at::timestamptz(3) <@ ixmis.validity
  $BODY$
  LANGUAGE sql IMMUTABLE;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'view', 20191014001, 'adopt deployment group renaming from schema ix');
COMMIT;
