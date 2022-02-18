# COMMAND

`asset::orchestration:property.set`

# REQUEST PATH

```
PUT /orchestration/:tomID/property/
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
{ "command":   "asset::orchestration:property.set",
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
