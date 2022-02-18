# COMMAND

`asset::container:property.remove`

# REQUEST PATH

`DELETE /container/:tomID/property/`

# REQUEST BODY

```
{ "container": {
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
{ "command":   "asset::container:property.remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
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
