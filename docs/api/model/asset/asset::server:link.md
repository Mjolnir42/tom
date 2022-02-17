# COMMAND

`asset::server:link`

# REQUEST PATH

```
POST /server/:tomID/link/
```

# REQUEST BODY

```
{ "server": {
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
{ "command":   "asset::server:link",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
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
