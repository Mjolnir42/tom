# COMMAND

`asset::runtime:stack`

# REQUEST PATH

```
PUT /runtime/:tomID/parent
```

# REQUEST BODY

```
{ "runtime": {
    "namespace": "<string>",
    "name":      "<string>",
    "property": {
      "<string>": {
        "attribute":  "asset::meta-cmd::stack",
        "value":      "<stackTargetTomID>",
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
{ "command":   "asset::runtime:stack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
    "namespace": "<string>",
    "name":      "<string>",
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
