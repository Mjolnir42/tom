--
--
-- VIEW SCHEMA
CREATE  VIEW view.information_system_details AS
SELECT  ixis.isID AS informationSystemID,
        md.dictionaryID AS dictionaryID,
        md.name AS dictionaryName,
        msa.attributeID AS attributeID,
        msa.attribute AS attributeName,
        ixiv.value AS attributeValue
FROM    yp.information_system AS ixis
JOIN    meta.dictionary AS md
  ON    ixis.dictionaryID = md.dictionaryID
JOIN    yp.information_system_standard_attribute_values AS ixiv
  ON    ixis.isID = ixiv.isID
JOIN    meta.standard_attribute AS msa
  ON    ixiv.attributeID = msa.attributeID
WHERE   NOW()::timestamptz(3) <@ ixiv.validity
UNION
SELECT  ixis.isID AS informationSystemID,
        md.dictionaryID AS dictionaryID,
        md.name AS dictionaryName,
        mqa.attributeID AS attributeID,
        mqa.attribute AS attributeName,
        ixqv.value AS attributeValue
FROM    yp.information_system AS ixis
JOIN    meta.dictionary AS md
  ON    ixis.dictionaryID = md.dictionaryID
JOIN    yp.information_system_unique_attribute_values AS ixqv
  ON    ixis.isID = ixqv.isID
JOIN    meta.unique_attribute AS mqa
  ON    ixqv.attributeID = mqa.attributeID
WHERE   NOW()::timestamptz(3) <@ ixqv.validity;

CREATE  FUNCTION view.information_system_details_at(at timestamptz)
  RETURNS TABLE ( informationSystemID uuid,
                  dictionaryID        uuid,
                  dictionaryName      text,
                  attributeID         uuid,
                  attributeName       text,
                  attributeValue      text)
  AS
  $BODY$
  SELECT  ixis.isID AS informationSystemID,
          md.dictionaryID AS dictionaryID,
          md.name AS dictionaryName,
          msa.attributeID AS attributeID,
          msa.attribute AS attributeName,
          ixiv.value AS attributeValue
  FROM    yp.information_system AS ixis
  JOIN    meta.dictionary AS md
    ON    ixis.dictionaryID = md.dictionaryID
  JOIN    yp.information_system_standard_attribute_values AS ixiv
    ON    ixis.isID = ixiv.isID
  JOIN    meta.standard_attribute AS msa
    ON    ixiv.attributeID = msa.attributeID
  WHERE   at::timestamptz(3) <@ ixiv.validity
  UNION
  SELECT  ixis.isID AS informationSystemID,
          md.dictionaryID AS dictionaryID,
          md.name AS dictionaryName,
          mqa.attributeID AS attributeID,
          mqa.attribute AS attributeName,
          ixqv.value AS attributeValue
  FROM    yp.information_system AS ixis
  JOIN    meta.dictionary AS md
    ON    ixis.dictionaryID = md.dictionaryID
  JOIN    yp.information_system_unique_attribute_values AS ixqv
    ON    ixis.isID = ixqv.isID
  JOIN    meta.unique_attribute AS mqa
    ON    ixqv.attributeID = mqa.attributeID
  WHERE   at::timestamptz(3) <@ ixqv.validity
  $BODY$
  LANGUAGE sql IMMUTABLE;
