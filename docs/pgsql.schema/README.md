# DATABASE SCHEMA NOTES

## Timezone

The database schema has check constraints to force persisted timestamps
to be in timezone UTC. The client session timezone setting can be set
via `SET TIME ZONE 'UTC';`.

The default database timezone can be configured in `postgresql.conf`.

All default value timestamps have been updated to be in time zone 'UTC'.

## INSTALLATION

### Roles and database creation

The `0x00.setup.sql` file contains the installation instructions
required to be performed by a superuser pgSQL role. The default setup
assumes that there are two roles for Tom, one owns the database and can
perform all the DDL commands. The other is the service account and only
has permissions to perform DML statements.

This file also contains the schema creation using the owner account.

### Table creation

After initial setup, the `0x*.sql - 0xff*.sql` files contain the
database schema. These files can be installed in lexicographic order
using for example the `\i /path/to/file/0x01.inventory.sql` command.
