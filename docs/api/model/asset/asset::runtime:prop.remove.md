# COMMAND

`asset::runtime:property.remove`

# REQUEST PATH

```
DELETE /runtime/:tomID/property/
```

# REQUEST BODY

```
{ "runtime": {
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
{ "command":   "asset::runtime:property.remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
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
