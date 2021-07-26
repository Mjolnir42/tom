--
--
-- DATABASE SETUP --
CREATE ROLE tom_owner WITH LOGIN PASSWORD 'xxx';
CREATE ROLE tomsvc    WITH LOGIN PASSWORD 'xyz';
CREATE DATABASE tom WITH OWNER tom_owner ENCODING UTF8 LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8' TEMPLATE template0;
\connect tom
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- update pg_hba.conf as required:
-- local   tom  tom_owner              password
-- host    tom  tomsvc      samehost   password
SELECT pg_reload_conf();

GRANT CONNECT ON DATABASE tom TO tom_owner;
GRANT CONNECT ON DATABASE tom TO tomsvc;
\connect tom tom_owner

-- create required function to index on uuid columns
CREATE OR REPLACE FUNCTION uuid_to_bytea(_uuid uuid) 
  RETURNS bytea AS                
  $BODY$
  select decode(replace(_uuid::text, '-', ''), 'hex');
  $BODY$
  LANGUAGE sql IMMUTABLE;

CREATE SCHEMA IF NOT EXISTS asset;
CREATE SCHEMA IF NOT EXISTS bulk;
CREATE SCHEMA IF NOT EXISTS filter;
CREATE SCHEMA IF NOT EXISTS inventory;
CREATE SCHEMA IF NOT EXISTS ix;
CREATE SCHEMA IF NOT EXISTS meta;
CREATE SCHEMA IF NOT EXISTS view;
CREATE SCHEMA IF NOT EXISTS yp;

SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory;
ALTER DATABASE tom SET search_path TO ix, meta, filter, yp, asset, 'view', bulk, inventory;
-- configure client session
SET TIME ZONE 'UTC';
