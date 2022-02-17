# COMMAND

`asset::server:property.remove`

# REQUEST PATH

```
DELETE /server/:tomID/property/
```

# REQUEST BODY

```
{ "server": {
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
      "<string>": {
        "attribute":  "<string>"
      },
      ...
    }
  }
}
```

# RESPONSE

```
{ "command":   "asset::server:property.remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
      "<string>": {
        "attribute": "<string>"
      },
      ...
    }
  }]
}
```
