BEGIN;
  ALTER TABLE meta.dictionary_unique_attribute_values DROP CONSTRAINT __mdq_temp_uniq;
  ALTER TABLE meta.dictionary_unique_attribute_values ADD CONSTRAINT  __mdq_temp_uniq
              EXCLUDE USING gist (public.uuid_to_bytea(dictionaryID) WITH =,
                                  public.uuid_to_bytea(attributeID) WITH =,
                                  value WITH =,
                                  validity WITH &&);

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'meta', 20191022001, 'fix constraint __mdq_temp_uniq');
COMMIT;
