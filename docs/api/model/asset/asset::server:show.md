# COMMAND

`asset::server:show`

# REQUEST PATH

```
GET /server/:tomID
GET /server/?namespace=<string>&name=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::server:show",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
    "parent":    "<string>",
    "createdAt": "<timespec>",
    "createdBy": "<string>",
    "resources": [
      "<string>",
      ...
    ],
    "children": [
      "<string>",
      ...
    ],
    "link": [
      "<string>",
      ...
    ],
    "property": {
      "<string>": {
        "attribute":   "<string>",
        "value":       "<string>",
        "namespace":   "<string>",
        "createdAt":   "<timespec>",
        "createdBy":   "<string>",
        "validSince":  "<timespec>",
        "validUntil": "<timespec>"
      },
      "<string>_list": {
        "attribute":       "<string>",
        "value":           "<string>",
        "structuredValue": [
          ...
        ],
        "namespace":       "<string>",
        "createdAt":       "<timespec>",
        "createdBy":       "<string>",
        "validSince":      "<timespec>",
        "validUntil":      "<timespec>"
      },
      "<string>_json": {
        "attribute":       "<string>",
        "value":           "<string>",
        "structuredValue": {
          ...
        },
        "namespace":       "<string>",
        "createdAt":       "<timespec>",
        "createdBy":       "<string>",
        "validSince":      "<timespec>",
        "validUntil":      "<timespec>"
      },
      ...
    }
  }]
}
```
