# DATABASE SCHEMA NOTES

## Timezone

The database schema has check constraints to force persisted timestamps
to be in timezone UTC. The client session timezone setting can be set
via `SET TIME ZONE 'UTC';`.
