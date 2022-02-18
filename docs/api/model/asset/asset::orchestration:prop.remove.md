# COMMAND

`asset::orchestration:property.remove`

# REQUEST PATH

```
DELETE /orchestration/:tomID/property/
```

# REQUEST BODY

```
{ "orchestration": {
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
{ "command":   "asset::orchestration:property.remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
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
