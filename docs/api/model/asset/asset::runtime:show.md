# COMMAND

`asset::runtime:show`

# REQUEST PATH

```
GET /runtime/:tomID
GET /runtime/?namespace=<string>&name=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::runtime:show",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
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
        ""validUntil": "<timespec>"
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
