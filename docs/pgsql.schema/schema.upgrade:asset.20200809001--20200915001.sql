BEGIN;
  ALTER TABLE asset.socket_parent DROP CONSTRAINT __asop_uq_parent;
  ALTER TABLE asset.socket_parent DROP CONSTRAINT __fk_asop_orchID;
  ALTER TABLE asset.socket_parent DROP COLUMN parentOrchestrationID;
  ALTER TABLE asset.socket_parent ALTER COLUMN parentRuntimeID SET NOT NULL;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'asset', 20200915001, 'remove orchestration environments as possible socket parent');
COMMIT;
