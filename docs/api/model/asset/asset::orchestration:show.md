# COMMAND

`asset::orchestration:show`

# REQUEST PATH

```
GET /orchestration/:tomID
GET /orchestration/?namespace=<string>&name=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::orchestration:show",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
    "createdAt": "<timespec>",
    "createdBy": "<string>",
    "parent": [
      "<string>",
      ...
    ],
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
        "validUntil":  "<timespec>"
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
