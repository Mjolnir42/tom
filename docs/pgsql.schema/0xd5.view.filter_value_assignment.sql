--
--
-- VIEW SCHEMA
-- -- filter_value_assignment lists all filter values and the entities
-- -- they are currently assigned to
CREATE  VIEW view.filter_value_assignment AS
SELECT  fvao.dictionaryID AS dictionaryID,
        fvao.filterID AS filterID,
        fvao.filterValueID AS filterValueID,
        fvao.tlsID AS tlsID,
        fvao.productID AS productID,
        fvao.isID AS isID,
        fvao.componentID AS componentID,
        fvao.groupID AS groupID,
        fvao.orchID AS orchID,
        fvao.rteID AS rteID,
        fvao.serverID AS serverID,
        fvao.validity AS validity
FROM    filter.value_assignment__one AS fvao
WHERE   NOW()::timestamptz(3) <@ fvao.validity
UNION
SELECT  fvam.dictionaryID AS dictionaryID,
        fvam.filterID AS filterID,
        fvam.filterValueID AS filterValueID,
        fvam.tlsID AS tlsID,
        fvam.productID AS productID,
        fvam.isID AS isID,
        fvam.componentID AS componentID,
        fvam.groupID AS groupID,
        fvam.orchID AS orchID,
        fvam.rteID AS rteID,
        fvam.serverID AS serverID,
        fvam.validity AS validity
FROM    filter.value_assignment__many AS fvam
WHERE   NOW()::timestamptz(3) <@ fvam.validity;

-- -- filter_value_assignment_at lists all filter values and the entities
-- -- they are assigned to at a specific point in time
CREATE  FUNCTION view.filter_value_assignment_at(at timestamptz)
  RETURNS TABLE ( dictionaryID  uuid,
                  filterID      uuid,
                  filterValueID uuid,
                  tlsID         uuid,
                  productID     uuid,
                  isID          uuid,
                  componentID   uuid,
                  groupID       uuid,
                  orchID        uuid,
                  rteID         uuid,
                  serverID      uuid,
                  validity      tstzrange)
  AS
  $BODY$
  SELECT  fvao.dictionaryID AS dictionaryID,
          fvao.filterID AS filterID,
          fvao.filterValueID AS filterValueID,
          fvao.tlsID AS tlsID,
          fvao.productID AS productID,
          fvao.isID AS isID,
          fvao.componentID AS componentID,
          fvao.groupID AS groupID,
          fvao.orchID AS orchID,
          fvao.rteID AS rteID,
          fvao.serverID AS serverID,
          fvao.validity AS validity
  FROM    filter.value_assignment__one AS fvao
  WHERE   at::timestamptz(3) <@ fvao.validity
  UNION
  SELECT  fvam.dictionaryID AS dictionaryID,
          fvam.filterID AS filterID,
          fvam.filterValueID AS filterValueID,
          fvam.tlsID AS tlsID,
          fvam.productID AS productID,
          fvam.isID AS isID,
          fvam.componentID AS componentID,
          fvam.groupID AS groupID,
          fvam.orchID AS orchID,
          fvam.rteID AS rteID,
          fvam.serverID AS serverID,
          fvam.validity AS validity
  FROM    filter.value_assignment__many AS fvam
  WHERE   at::timestamptz(3) <@ fvam.validity;
  $BODY$
  LANGUAGE sql IMMUTABLE;
