# COMMAND

`asset::server:stack`

# REQUEST PATH

```
PUT /server/:tomID/parent
```

# REQUEST BODY

```
{ "server": {
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
{ "command":   "asset::server:stack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
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
