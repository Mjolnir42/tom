# COMMAND

`asset::orchestration:list`

# REQUEST PATH

```
GET /orchestration/
GET /orchestration/?namespace=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::orchestration:list",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration-list": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
    "createdAt": "<timespec>",
    "createdBy": "<string>"
  },
  ...
  ]
}
```
