# slamdd-go

slamdd-go is the server data daemon. It also has facilities for receiving,
filtering, processing and sending of IPFIX data.

## Configuration

Example configuration file is provided in `docs/slamdd/configuration/slam.conf.example`.

### Data daemon integration

Unless are TOM service is running, disable the data daemon functionality.

```
data.daemon.enabled: false
api: https://127.0.0.1:8443/
api.ca.file: /svc/slamdd/rootCA.pem
log.level: info
log.path: /var/log/slamdd
authentication: {
    credential.path: /svc/slamdd
}
```

### IPFIX configuration

If enabled, all of the enabled server processes will be started. At most one
server per listing protocol.

If forwarding is enabled, the configured enabled clients will be started.

If filter rules are added, the filter rules are applied to the incoming
IPFIX messages.

```
ipfix: {
  enabled: true
  forwarding.enabled: true
  processing.enabled: true
  # filter, aggregate, filter+aggregate
  processing.type: filter
  template.file: /svc/slamdd/ipfix.tmpl
  server: [
  ]
  client: [
  ]
  filter: {
    rules: [
    ]
  }
```

### IPFIX/udp server

```
{ enabled: true
  listen.protocol: udp
  listen.address: 127.0.0.1:4739
}
```

### IPFIX/tcp server

```
{ enabled: true
  listen.protocol: tcp
  listen.address: 127.0.0.1:4739
}
```

### IPFIX/tls server

```
{ enabled: true
  listen.protocol: tls
  listen.address: 127.0.0.1:4740
  tls.servername: localhost
  ca.file: /svc/slamdd/rootCA.pem
  certificate.file: /svc/slamdd/localhost.pem
  certificate.keyfile: /svc/slamdd/localhost.key
}
```

### IPFIX/json server

This server accepts PUT, POST or PATCH requests on
`https://127.0.0.1:8443/submit`. It will not start without configured basic
auth and all clients need to use it.
Data received on the IPFIX/json input will be copied directly and only
to the IPFIX/json output client. No filtering, processing, re-encoding or
other functions are applied to it.

```
{ enabled: true
  listen.protocol: json
  listen.address: 127.0.0.1:8443
  tls.servername: localhost
  ca.file: /svc/slamdd/rootCA.pem
  certificate.file: /svc/slamdd/localhost.pem
  certificate.keyfile: /svc/slamdd/localhost.key
  basic.auth.user: username
  basic.auth.pass: secret_password
}
```

#### IPFIX/udp client

```
{ enabled: true
  forwarding.protocol: udp
  forwarding.address: 127.0.0.1:4739
  unfiltered.copy: false
}
```

#### IPFIX/tcp client

```
{ enabled: true
  forwarding.protocol: tcp
  forwarding.address: 127.0.0.1:4739
  unfiltered.copy: false
}
```

#### IPFIX/tls client

```
{ enabled: true
  forwarding.protocol: tls
  forwarding.address: 127.0.0.1:4740
  ca.file: /svc/slamdd/rootCA.pem
  unfiltered.copy: false
}
```

#### IPFIX/json client

This client is intended to connect to the `http_input` plugin from
logstash/elastic.

It will send the JSON body via POST request to
`https://127.0.0.1:8443/request/path`.

```
{ enabled: true
  forwarding.protocol: json
  forwarding.address: 127.0.0.1:8443
  ca.file: /svc/slamdd/rootCA.pem
  json.format: flowdata
  basic.auth.user: username
  basic.auth.pass: secret_password
  endpoint: /request/path
  method: POST
}
```
