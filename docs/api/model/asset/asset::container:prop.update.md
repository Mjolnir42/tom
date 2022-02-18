# COMMAND

`asset::container:property.update`

# REQUEST PATH

```
PATCH /container/:tomID/property/
```

# REQUEST BODY

```
{ "container": {
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
{ "command":   "asset::container:property.update",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
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
