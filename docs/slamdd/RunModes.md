# Runtime Configuration Modes for slamdd-go

## IPFIX

```
ipfix: {
    enabled:                  yes|no
    forwarding.enabled:       yes|no
    processing.enabled:       yes|no
    processing.type:          filter
    template.file:            /path/...
    server: [
    {   enabled:              yes|no
        listen.protocol:      udp
        listen.address:       127.0.0.1:4739
    },
    {   enabled:              yes|no
        listen.protocol:      tcp
        listen.address:       127.0.0.1:4739
    },
    {   enabled:              yes|no
        listen.protocol:      tls
        listen.address:       127.0.0.1:4740
        tls.servername:       localhost
        ca.file:              /path/...
        certificate.file:     /path/...
        certificate.keyfile:  /path/...
    }]
    client: [
    {   enabled:              yes|no
        forwarding.address:   ...
        forwarding.protocol:  udp
        unfiltered.copy:      yes|no
    },
    {   enabled:              yes|no
        forwarding.address:   ...
        forwarding.protocol:  tcp
        unfiltered.copy:      yes|no
    },
    {   enabled:              yes|no
        forwarding.address:   ...
        forwarding.protocol:  tls
        ca.file:              /path/...
        unfiltered.copy:      yes|no
    },
    {   enabled:              yes|no
        forwarding.address:   ...
        forwarding.protocol:  json
        ca.file:              /path/...
        unfiltered.copy:      yes|no
        json.format:          vflow|flowdata
    }]
    filter: {
      rules: [
      ]
    }
}
```

## Forwarding IPFIX, no filtering, no JSON

```
  enabled: yes
  forwarding.enabled: yes
  processing.enabled: false
        -> filtering: false
      -> aggregation: false
  server.udp:   yes|no
  server.tcp:   yes|no
  server.tls:   yes|no
  server.json:  no
  client.udp:   yes|no
  client.tcp:   yes|no
  client.tls:   yes|no
  client.json:  no

              | Reader        | Writer        |
--------------+---------------+---------------+
chan  inUDP   | mux           | udp.server    |
chan  inTCP   | mux           | tcp.server    |
chan  inTLS   | mux           | tls.server    |
chan  inJSN   | connectJSON   | json.server   |
              |               |               |
chan  outUDP  | udp.client    | mux           |
              | opportunUDP   |               |
chan  outTCP  | tcp.client    | mux           |
              | opportunTCP   |               |
chan  outTLS  | tls.client    | mux           |
              | opportunTLS   |               |
chan  outJSN  | json.client   | connectJSON   |
              | opportunJSON  |               |
              |               |               |
chan  outFLT  | connectFilter | mux           |
chan  inFLT   | mux           | connectFilter |
chan  inFLJ   | mux           | -             |
chan  inFLR   | mux           | -             |
              |               |               |
chan  outAGG  | opportunAGG   | -             |
```

## Forwarding IPFIX, no filtering, JSON input

```
  enabled: yes
  forwarding.enabled: yes
  processing.enabled: false
        -> filtering: false
      -> aggregation: false
  server.udp:   yes|no
  server.tcp:   yes|no
  server.tls:   yes|no
  server.json:  yes
  client.udp:   yes|no
  client.tcp:   yes|no
  client.tls:   yes|no
  client.json:  no

Invalid configuration error.
```

## Forwarding IPFIX, no filtering, JSON output

```
  enabled: yes
  forwarding.enabled: yes
  processing.enabled: false
        -> filtering: false
      -> aggregation: false
  server.udp:   yes|no
  server.tcp:   yes|no
  server.tls:   yes|no
  server.json:  no
  client.udp:   yes|no
  client.tcp:   yes|no
  client.tls:   yes|no
  client.json:  yes

              | Reader        | Writer        |
--------------+---------------+---------------+
chan  inUDP   | mux           | udp.server    |
chan  inTCP   | mux           | tcp.server    |
chan  inTLS   | mux           | tls.server    |
chan  inJSN   | connectJSON   | json.server   |
              |               |               |
chan  outUDP  | udp.client    | mux           |
              | opportunUDP   |               |
chan  outTCP  | tcp.client    | mux           |
              | opportunTCP   |               |
chan  outTLS  | tls.client    | mux           |
              | opportunTLS   |               |
chan  outJSN  | json.client   | connectJSON   |
              | opportunJSON  |               |
              |               |               |
chan  outFLT  | filter        | mux           |
chan  inFLT   | mux           | filter        |
chan  inFLJ   | mux           | filter        |
chan  inFLR   | mux           | filter        |
              |               |               |
chan  outAGG  | opportunAGG   | -             |
```

