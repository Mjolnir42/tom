BEGIN;
  CREATE TABLE IF NOT EXISTS meta.attribute (
      dictionaryID                  uuid        NOT NULL,
      attribute                     text        NOT NULL,
      CONSTRAINT __uniq_attr_name   UNIQUE      ( dictionaryID, attribute )
  );
  ALTER TABLE meta.standard_attribute ADD CONSTRAINT __fk_msa_attr
              FOREIGN KEY ( dictionaryID, attribute )
              REFERENCES meta.attribute ( dictionaryID, attribute )
              ON DELETE CASCADE;
  ALTER TABLE meta.unique_attribute ADD CONSTRAINT __fk_msqa_attr
              FOREIGN KEY ( dictionaryID, attribute )
              REFERENCES meta.attribute ( dictionaryID, attribute )
              ON DELETE CASCADE;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'meta', 20201210001, 'add meta.attribute registry table');
COMMIT;
