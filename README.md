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

### METADATA SERVICE SETUP

Use `docs/tomd/configuration/tomd.conf.example` as sample service configuration.
Start the service using `tomd -c tomd.conf`.

More details in `docs/tom/README.md`.

### CLI SETUP

Use `docs/tom/configuration/tom.conf.example` as sample CLI
configuration. Create a `~/.tom/tom.conf` from it.

More details in `docs/tom/README.md`.

### DATA DAEMON SETUP

Use `docs/slamdd/configuration/slam.conf.example` as sample data daemon
configuration.

More details in `docs/slamdd/README.md`.

#### ENABLE zsh AUTOCOMPLETION

```
PROG=tom
_CLI_ZSH_AUTOCOMPLETE_HACK=1
source <(tom output-autocomplete)
alias tom='noglob tom'
```

## USER MANAGEMENT

Usermanagement is partially hooked up yet, machine self-signup works,
but otherwise only the `root` user is usable.
