BEGIN;
  ALTER TABLE yp.service_linking ADD COLUMN validity tstzrange NOT NULL DEFAULT tstzrange(NOW()::timestamptz(3), 'infinity', '[]');
  ALTER TABLE yp.service_linking ADD CONSTRAINT __validFrom_utc CHECK ( EXTRACT( TIMEZONE FROM lower( validity ) ) = '0' );
  ALTER TABLE yp.service_linking ADD CONSTRAINT __validUntil_utc CHECK ( EXTRACT( TIMEZONE FROM upper( validity ) ) = '0' );
  ALTER TABLE yp.service_linking ADD CONSTRAINT __ypsl_temporal EXCLUDE USING gist (public.uuid_to_bytea(serviceID) WITH =, public.uuid_to_bytea(endpointID) WITH =, validity WITH &&);
  ALTER TABLE yp.service_linking DROP CONSTRAINT __pk_ypsl;
  ALTER TABLE yp.service_linking RENAME CONSTRAINT __fk_ypsl_serID TO __fk_ypsm_serID;
  ALTER TABLE yp.service_linking RENAME CONSTRAINT __fk_ypsl_enpID TO __fk_ypsm_enpID;
  ALTER TABLE yp.service_linking RENAME CONSTRAINT __ypsl_temporal TO __ypsm_temporal;
  ALTER TABLE yp.service_linking RENAME TO service_mapping;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'yp', 20201016001, 'switch from service linking 1:n to mapping n:m');
COMMIT;
