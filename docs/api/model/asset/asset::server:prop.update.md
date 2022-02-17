# COMMAND

`asset::server:property.update`

# REQUEST PATH

```
PATCH /server/:tomID/property/
```

# REQUEST BODY

```
{ "server": {
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
{ "command":   "asset::server:property.update",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
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
