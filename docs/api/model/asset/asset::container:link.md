# COMMAND

`asset::container:link`

# REQUEST PATH

```
POST /container/:tomID/link/
```

# REQUEST BODY

```
{ "container": {
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
      "<string>": {
        "attribute":  "asset::meta-cmd::link",
        "value":      "<linkTargetTomID>",
        "validSince": "perpetual"
        "validUntil": "perpetual"
      },
      ...
    }
  }
}
```

# RESPONSE

```
{ "command":   "asset::container:link",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
      "<string>": {
        "attribute":  "asset::meta-cmd::link",
        "value":      "<linkTargetTomID>",
        "validSince": "perpetual"
        "validUntil": "perpetual"
      },
      ...
    }
  }]
}
```
