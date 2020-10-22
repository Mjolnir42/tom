BEGIN;
  CREATE TABLE IF NOT EXISTS asset.server_parent (
      serverID                      uuid        NOT NULL,
      parentRuntimeID               uuid        NULL,
      validity                      tstzrange   NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]'),
      CONSTRAINT __fk_asp_srvID     FOREIGN KEY ( serverID ) REFERENCES asset.server ( serverID ) ON DELETE RESTRICT,
      CONSTRAINT __fk_asp_rtenv     FOREIGN KEY ( parentRuntimeID ) REFERENCES asset.runtime_environment ( rteID ) ON DELETE RESTRICT,
      CONSTRAINT __validFrom_utc    CHECK       ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' ),
      CONSTRAINT __validUntil_utc   CHECK       ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' ),
      CONSTRAINT __asp_uq_parent    CHECK       ( parentRuntimeID IS NOT NULL ),
      CONSTRAINT __asp_temporal     EXCLUDE     USING gist (public.uuid_to_bytea(serverID) WITH =,
                                                            validity WITH &&)
  );

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'asset', 20191104001, 'add server_parent table');
COMMIT;
