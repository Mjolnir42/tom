# COMMAND

`asset::runtime:add`

# REQUEST PATH

```
POST /runtime/
```

# REQUEST BODY

```
{ "runtime": {
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
{ "command":   "asset::runtime:add",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
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
  }]
}
```
