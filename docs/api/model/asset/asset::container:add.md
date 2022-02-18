# COMMAND

`asset::container:add`

# REQUEST PATH

```
POST /container/
```

# REQUEST BODY

```
{ "container": {
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
{ "command":   "asset::container:add",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
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
