BEGIN;
  ALTER TABLE filter.confidentiality_mapping RENAME COLUMN subgroupID TO groupID;
  ALTER TABLE filter.integrity_mapping.sql RENAME COLUMN subgroupID TO groupID;
  ALTER TABLE filter.availability_mapping RENAME COLUMN subgroupID TO groupID;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'filter', 20191014001, 'adopt deployment group renaming from schema ix');
COMMIT;
