# COMMAND

`asset::orchestration:stack`

# REQUEST PATH

```
PUT /orchestration/:tomID/parent
```

# REQUEST BODY

```
{ "orchestration": {
    "namespace": "<string>",
    "name": "<string>",
    "property": {
      "<string>": {
        "attribute": "asset::meta-cmd::stack",
        "value": "<stackTargetTomID>",
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
{ "command":   "asset::orchestration:stack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name": "<string>",
    "property": {
      "<string>": {
        "attribute":  "asset::meta-cmd::stack",
        "value":      "<stackTargetTomID>",
        "validSince": "<timespec,optional>"
        "validUntil": "<timespec,optional>"
      },
      ...
    }
  }]
}
```
