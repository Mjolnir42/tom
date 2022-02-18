# COMMAND

`asset::container:stack`

# REQUEST PATH

```
PUT /container/:tomID/parent
```

# REQUEST BODY

```
{ "container": {
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
{ "command":   "asset::container:stack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
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
