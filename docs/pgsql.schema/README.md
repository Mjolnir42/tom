# DATABASE SCHEMA NOTES

## Timezone

The database schema has check constraints to force persisted timestamps
to be in timezone UTC. The client session timezone setting can be set
via `SET TIME ZONE 'UTC';`.

The default database timezone can be configured in `postgresql.conf`.

All default value timestamps have been updated to be in time zone 'UTC'.

## INSTALLATION

### Roles and database creation

The `0x00.setup.sql` file contains installation instructions
required to be performed by a superuser pgSQL role. The default setup
assumes that there are two roles for Tom, one owns the database and can
perform all the DDL commands. The other is the service account and only
has permissions to perform DML statements.

This file also contains the schema creation using the owner account.

These roles are called `tom_owner` and `tomsvc` within the scripts and
grant statements.
All instructions in `0x00.setup.sql` after the `\connect tom tom_owner`
as well as all subsequent SQL files are assumed to be executed by the
owner role, and not a superuser account.

### Table creation

After initial setup, the `0x*.sql - 0xff*.sql` files contain the
database schema. These files can be installed in lexicographic order
using for example the `\i /path/to/file/0x01.inventory.sql` command.

### Schema Upgrade

The `schema.upgrade:*.sql` files are SQL scripts for upgrading an existing
database setup. During a clean install, they do not neet to be executed.
