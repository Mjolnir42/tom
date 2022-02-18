# COMMAND

`asset::orchestration:add`

# REQUEST PATH

```
POST /orchestration/
```

# REQUEST BODY

```
{ "orchestration": {
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
  }
}
```

# RESPONSE

```
{ "command":   "asset::orchestration:add",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
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
