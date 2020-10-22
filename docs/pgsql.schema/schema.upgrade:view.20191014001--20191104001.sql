BEGIN;
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

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'view', 20191104001, 'add dictionary schema/definition functions');
COMMIT;
