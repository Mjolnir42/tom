--
--
-- VIEW SCHEMA
-- -- dictionary_details lists all for all currently defined
-- -- dictionaries all self-describing attributes (dict_*) with
-- -- their current values.
CREATE  VIEW view.dictionary_details AS
SELECT  md.dictionaryID AS dictionaryID,
        msa.attributeID AS attributeID,
        msa.attribute AS attributeName,
        mda.value AS attributeValue
FROM    meta.dictionary AS md
JOIN    meta.dictionary_standard_attribute_values AS mda
  ON    md.dictionaryID = mda.dictionaryID
JOIN    meta.standard_attribute AS msa
  ON    mda.attributeID = msa.attributeID
WHERE   NOW()::timestamptz(3) <@ mda.validity
UNION
SELECT  md.dictionaryID AS dictionaryID,
        mqa.attributeID AS attributeID,
        mqa.attribute AS attributeName,
        mdq.value AS attributeValue
FROM    meta.dictionary AS md
JOIN    meta.dictionary_unique_attribute_values AS mdq
  ON    md.dictionaryID = mdq.dictionaryID
JOIN    meta.unique_attribute AS mqa
  ON    mdq.attributeID = mqa.attributeID
WHERE   NOW()::timestamptz(3) <@ mdq.validity;

-- -- dictionary_details_at lists all for all currently defined
-- -- dictionaries all self-describing attributes (dict_*) with
-- -- their values at a specified past timestamp.
CREATE  FUNCTION view.dictionary_details_at(at timestamptz)
  RETURNS TABLE ( dictionaryID    uuid,
                  attributeID     uuid,
                  attributeName   text,
                  attributeValue  text)
  AS
  $BODY$
  SELECT  md.dictionaryID AS dictionaryID,
          msa.attributeID AS attributeID,
          msa.attribute AS attributeName,
          mda.value AS attributeValue
  FROM    meta.dictionary AS md
  JOIN    meta.dictionary_standard_attribute_values AS mda
    ON    md.dictionaryID = mda.dictionaryID
  JOIN    meta.standard_attribute AS msa
    ON    mda.attributeID = msa.attributeID
  WHERE   at::timestamptz(3) <@ mda.validity
  UNION
  SELECT  md.dictionaryID AS dictionaryID,
          mqa.attributeID AS attributeID,
          mqa.attribute AS attributeName,
          mdq.value AS attributeValue
  FROM    meta.dictionary AS md
  JOIN    meta.dictionary_unique_attribute_values AS mdq
    ON    md.dictionaryID = mdq.dictionaryID
  JOIN    meta.unique_attribute AS mqa
    ON    mdq.attributeID = mqa.attributeID
  WHERE   at::timestamptz(3) <@ mdq.validity;
  $BODY$
  LANGUAGE sql IMMUTABLE;

-- -- dictionary_schema_of lists the attribute schema of
-- -- a dictionary by name.
-- -- It returns the attribute names and their type.
CREATE  FUNCTION view.dictionary_schema_of(arg text)
  RETURNS TABLE ( dictionaryID    uuid,
                  dictionaryName  text,
                  attributeID     uuid,
                  attributeName   text,
                  attributeType  text)
  AS
  $BODY$
  SELECT  md.dictionaryID  AS dictionaryID,
          md.name          AS dictionaryName,
          msa.attributeID  AS attributeID,
          msa.attribute    AS attributeName,
          'standard'::text AS attributeType
  FROM    meta.dictionary  AS md
    JOIN  meta.standard_attribute AS msa
      ON  md.dictionaryID = msa.dictionaryID
  WHERE   msa.attribute NOT LIKE 'dict_%'
    AND   md.name = arg
  UNION
  SELECT  md.dictionaryID  AS dictionaryID,
          md.name          AS dictionaryName,
          mqa.attributeid  AS attributeID,
          mqa.attribute    AS attributename,
          'unique'::text   AS attributeType
  FROM    meta.dictionary  AS md
    JOIN  meta.unique_attribute AS mqa
      ON  md.dictionaryID = mqa.dictionaryID
    WHERE mqa.attribute NOT LIKE 'dict_%'
    AND   md.name = arg;
  $BODY$
  LANGUAGE sql IMMUTABLE;

-- -- dictionary_definition_of lists the dictionary definition
-- -- by name.
-- -- It returns all `dict_` attributes, their type as well as
-- -- their value if it is currently set.
CREATE  FUNCTION view.dictionary_definition_of(arg text)
  RETURNS TABLE ( dictionaryID    uuid,
                  dictionaryName  text,
                  attributeID     uuid,
                  attributeName   text,
                  attributeValue  text,
                  attributeType   text)
  AS
  $BODY$
  SELECT  md.dictionaryID  AS dictionaryID,
          md.name          AS dictionaryName,
          msa.attributeID  AS attributeID,
          msa.attribute    AS attributeName,
          mda.value        AS attributevalue,
          'standard'::text AS attributeType
  FROM    meta.dictionary  AS md
    JOIN  meta.standard_attribute AS msa
      ON  md.dictionaryID = msa.dictionaryID
    LEFT  JOIN meta.dictionary_standard_attribute_values AS mda
      ON  msa.attributeID = mda.attributeID
  WHERE   msa.attribute LIKE 'dict_%'
    AND   md.name = arg
    AND   ( NOW()::timestamptz(3) <@ mda.validity OR mda.value IS NULL )
  UNION
  SELECT  md.dictionaryID  AS dictionaryID,
          md.name          AS dictionaryName,
          mqa.attributeid  AS attributeID,
          mqa.attribute    AS attributename,
          mdq.value        AS attributevalue,
          'unique'::text   AS attributeType
  FROM    meta.dictionary  AS md
    JOIN  meta.unique_attribute AS mqa
      ON  md.dictionaryID = mqa.dictionaryID
    LEFT  JOIN meta.dictionary_unique_attribute_values AS mdq
      ON  mqa.attributeID = mdq.attributeID
    WHERE mqa.attribute LIKE 'dict_%'
    AND   md.name = arg
    AND   ( NOW()::timestamptz(3) <@ mdq.validity OR mdq.value IS NULL );
  $BODY$
  LANGUAGE sql IMMUTABLE;

-- -- dictionary_definition_of lists the dictionary definition
-- -- by name.
-- -- It returns all `dict_` attributes, their type as well as
-- -- their value if it is set at timestamp `at`.
CREATE  FUNCTION view.dictionary_definition_of_at(arg text, at timestamptz)
  RETURNS TABLE ( dictionaryID    uuid,
                  dictionaryName  text,
                  attributeID     uuid,
                  attributeName   text,
                  attributeValue  text,
                  attributeType   text)
  AS
  $BODY$
  SELECT  md.dictionaryID  AS dictionaryID,
          md.name          AS dictionaryName,
          msa.attributeID  AS attributeID,
          msa.attribute    AS attributeName,
          mda.value        AS attributevalue,
          'standard'::text AS attributeType
  FROM    meta.dictionary  AS md
    JOIN  meta.standard_attribute AS msa
      ON  md.dictionaryID = msa.dictionaryID
    LEFT  JOIN meta.dictionary_standard_attribute_values AS mda
      ON  msa.attributeID = mda.attributeID
  WHERE   msa.attribute LIKE 'dict_%'
    AND   md.name = arg
    AND   ( at::timestamptz(3) <@ mda.validity OR mda.value IS NULL )
  UNION
  SELECT  md.dictionaryID  AS dictionaryID,
          md.name          AS dictionaryName,
          mqa.attributeid  AS attributeID,
          mqa.attribute    AS attributename,
          mdq.value        AS attributevalue,
          'unique'::text   AS attributeType
  FROM    meta.dictionary  AS md
    JOIN  meta.unique_attribute AS mqa
      ON  md.dictionaryID = mqa.dictionaryID
    LEFT  JOIN meta.dictionary_unique_attribute_values AS mdq
      ON  mqa.attributeID = mdq.attributeID
    WHERE mqa.attribute LIKE 'dict_%'
    AND   md.name = arg
    AND   ( at::timestamptz(3) <@ mdq.validity OR mdq.value IS NULL );
  $BODY$
  LANGUAGE sql IMMUTABLE;
