# COMMAND

`asset::orchestration:property.update`

# REQUEST PATH

```
PATCH /orchestration/:tomID/property/
```

# REQUEST BODY

```
{ "orchestration": {
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
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
{ "command":   "asset::orchestration:property.update",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
      "<string>": {
        "attribute": "<string>",
        "value":     "<string>"
      },
      ...
    }
  }]
}
```
