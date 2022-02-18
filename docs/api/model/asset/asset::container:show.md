# COMMAND

`asset::container:show`

# REQUEST PATH

```
GET /container/:tomID
GET /container/?namespace=<string>&name=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::container:show",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
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
