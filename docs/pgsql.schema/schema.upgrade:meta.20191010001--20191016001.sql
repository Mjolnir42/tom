BEGIN;
  ALTER TABLE meta.dictionary_attribute RENAME TO dictionary_attr_values;
  ALTER TABLE meta.dictionary_unique RENAME TO dictionary_attr_uniq_values;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'meta', 20191016001, 'use better names for attribute value assignment tables');
COMMIT;
