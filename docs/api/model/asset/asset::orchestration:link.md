# COMMAND

`asset::orchestration:link`

# REQUEST PATH

```
POST /orchestration/:tomID/link/
```

# REQUEST BODY

```
{ "orchestration": {
    "namespace": "<string>",
    "name": "<string>",
    "property": {
      "<string>": {
        "attribute": "asset::meta-cmd::link",
        "value": "<linkTargetTomID>",
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
{ "command":   "asset::orchestration:link",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name": "<string>",
    "property": {
      "<string>": {
        "attribute": "asset::meta-cmd::link",
        "value": "<linkTargetTomID>",
        "validSince": "perpetual"
        "validUntil": "perpetual"
      },
      ...
    }
  }]
}
```
