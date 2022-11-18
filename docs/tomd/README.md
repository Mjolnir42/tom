# tomd

This document describes the installation of the `tomd` metadata service.

Inside the configuration file, lines starting with `#` are comments.

## Database setup

Install the database schema according to the instructions in
`docs/pgsql.schema/README.md`.

## Service configuration

An example service configuration file is provided in
`docs/tomd/configuration/tomd.conf.example`.

### database connection configuration

```
database {
  host: 127.0.0.1
  user: tomsvc
  database: tom
  port: 5432
  password: **********
  timeout: 1
  # require, verify-full, verify-ca, disable
  tlsmode: verify-full
}
```

### service listening configuration

The service can listen on multiple addresses and/or ports at the same time,
either with TLS or without. Provide one `daemon {}` section for every
listener.

```
daemon {
  listen: 127.0.0.1
  port: 80
  tls: false
}
daemon {
  listen: 127.0.0.1
  port: 443
  tls: true
  cert.file: /svc/tomd/tls-cert.pem
  key.file: /svc/tomd/tls-key.pem
}
```

### log configuration

The log level may be adjusted, but there are not that many log messages
configured for lower, noisy levels.
The logfiles are store in log.path and reopened on SIGUSR2 to allow for
logrotate.

```
# debug, info, warn, error, fatal, panic
log.level: info
log.path: /svc/tomd/log
```

### internal configuration

```
handler.queue.length: 1
enforcement: false
```

`handler.queue.length` regulates the length of internal request buffer
queues.

`enforcement` regulates if authentication is enforced. If authentication is
not enforced, only user/server enrolment functions work.

## First start

1. startup the service with `enforcement: false`
2. initialize your root account according to `docs/tom/README.md`
3. restart the service with `enforcement: true`
