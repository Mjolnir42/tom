# COMMAND

`asset::server:add`

# REQUEST PATH

```
POST /server/
```

# REQUEST BODY

```
{ "server": {
    "namespace": "<string>",
    "property": {
      "name": {
        "attribute":  "name",
        "value":      "<string>",
        "validSince": "<timespec,optional>"
        "validUntil": "<timespec,optional>"
      },
      "type": {
        "attribute":  "type",
        "value":      "<string>",
        "validSince": "perpetual",
        "validUntil": "perpetual",
      },
      "<string>": {
        "attribute":  "<string>",
        "value":      "<string>",
        "validSince": "<timespec,optional>"
        "validUntil": "<timespec,optional>"
      },
      ...
    }
  }
}
```

# RESPONSE

```
{ "command":   "asset::server:add",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
    "namespace": "<string>",
    "property": {
      "name": {
        "attribute": "name",
        "value": "<string>",
        "validSince": "<timespec,optional>"
        "validUntil": "<timespec,optional>"
      },
      "type": {
        "attribute": "type",
        "value": "<string>",
        "validSince": "perpetual",
        "validUntil": "perpetual",
      },
      "<string>": {
        "attribute": "<string>",
        "value": "<string>",
        "validSince": "<timespec,optional>"
        "validUntil": "<timespec,optional>"
      },
      ...
    }
  }]
}
```
