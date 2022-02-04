# TOM

Hello, I'm Tom.

## INSTALLATION

### COMPILATION

```
git pull
go mod download
make install_all
```

### DATABASE SETUP

The database setup is described in `docs/pgsql.schema/README.md`

### SERVICE SETUP

Use `docs/tomd/configuration/tomd.conf` as sample service configuration.
Start the service using `tomd -c tomd.conf`.

### CLI SETP

Use `docs/tom/configuration/tom.conf.example` as sample CLI
configuration. Create a `~/.tom/tom.conf` from it.

## USER MANAGEMENT

Usermanagement is not hooked up yet, everybody is user `nobody`.
