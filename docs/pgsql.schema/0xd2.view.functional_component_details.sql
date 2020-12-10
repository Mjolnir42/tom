--
--
-- VIEW SCHEMA
CREATE  VIEW view.functional_component_details AS
SELECT  ixfc.componentID AS componentID,
        md.dictionaryID AS componentDictionaryID,
        md.name AS componentDictionaryName,
        ypisp.isID AS informationSystemID,
        msa.attributeID AS attributeID,
        msa.attribute AS attributeName,
        ixfcqv.value AS attributeValue
FROM    ix.functional_component AS ixfc
JOIN    meta.dictionary AS md
  ON    ixfc.dictionaryID = md.dictionaryID
JOIN    ix.functional_component_unique_attribute_values AS ixfcqv
  ON    ixfc.componentID = ixfcqv.componentID
JOIN    meta.standard_attribute AS msa
  ON    ixfcqv.attributeID = msa.attributeID
JOIN    yp.information_system_parent AS ypisp
  ON    ixfc.componentID = ypisp.componentID
WHERE   NOW()::timestamptz(3) <@ ixfcqv.validity
  AND   NOW()::timestamptz(3) <@ ypisp.validity
UNION
SELECT  ixfc.componentID AS componentID,
        md.dictionaryID AS componentDictionaryID,
        md.name AS componentDictionaryName,
        ypisp.isID AS informationSystemID,
        mqa.attributeID AS attributeID,
        mqa.attribute AS attributeName,
        ixfcqv.value AS attributeValue
FROM    ix.functional_component AS ixfc
JOIN    meta.dictionary AS md
  ON    ixfc.dictionaryID = md.dictionaryID
JOIN    ix.functional_component_unique_attribute_values AS ixfcqv
  ON    ixfc.componentID = ixfcqv.componentID
JOIN    meta.unique_attribute AS mqa
  ON    ixfcqv.attributeID = mqa.attributeID
JOIN    yp.information_system_parent AS ypisp
  ON    ixfc.componentID = ypisp.componentID
WHERE   NOW()::timestamptz(3) <@ ixfcqv.validity
  AND   NOW()::timestamptz(3) <@ ypisp.validity;

CREATE  FUNCTION view.functional_component_details_at(at timestamptz)
  RETURNS TABLE ( componentID             uuid,
                  componentDictionaryID   uuid,
                  componentDictionaryName text,
                  informationSystemID     uuid,
                  attributeID             uuid,
                  attributeName           text,
                  attributeValue          text)
  AS
  $BODY$
  SELECT  ixfc.componentID AS componentID,
          md.dictionaryID AS componentDictionaryID,
          md.name AS componentDictionaryName,
          ypisp.isID AS informationSystemID,
          msa.attributeID AS attributeID,
          msa.attribute AS attributeName,
          ixfcqv.value AS attributeValue
  FROM    ix.functional_component AS ixfc
  JOIN    meta.dictionary AS md
    ON    ixfc.dictionaryID = md.dictionaryID
  JOIN    ix.functional_component_unique_attribute_values AS ixfcqv
    ON    ixfc.componentID = ixfcqv.componentID
  JOIN    meta.standard_attribute AS msa
    ON    ixfcqv.attributeID = msa.attributeID
  JOIN    yp.information_system_parent AS ypisp
    ON    ixfc.componentID = ypisp.componentID
  WHERE   at::timestamptz(3) <@ ixfcqv.validity
    AND   at::timestamptz(3) <@ ypisp.validity
  UNION
  SELECT  ixfc.componentID AS componentID,
          md.dictionaryID AS componentDictionaryID,
          md.name AS componentDictionaryName,
          ypisp.isID AS informationSystemID,
          mqa.attributeID AS attributeID,
          mqa.attribute AS attributeName,
          ixfcqv.value AS attributeValue
  FROM    ix.functional_component AS ixfc
  JOIN    meta.dictionary AS md
    ON    ixfc.dictionaryID = md.dictionaryID
  JOIN    ix.functional_component_unique_attribute_values AS ixfcqv
    ON    ixfc.componentID = ixfcqv.componentID
  JOIN    meta.unique_attribute AS mqa
    ON    ixfcqv.attributeID = mqa.attributeID
  JOIN    yp.information_system_parent AS ypisp
    ON    ixfc.componentID = ypisp.componentID
  WHERE   at::timestamptz(3) <@ ixfcqv.validity
    AND   at::timestamptz(3) <@ ypisp.validity
  $BODY$
  LANGUAGE sql IMMUTABLE;
