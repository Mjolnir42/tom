# COMMAND

`asset::runtime:link`

# REQUEST PATH

```
POST /runtime/:tomID/link/
```

# REQUEST BODY

```
{ "runtime": {
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
{ "command":   "asset::runtime:link",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
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
