# COMMAND

`asset::runtime:list`

# REQUEST PATH

```
GET /runtime/
GET /runtime/?namespace=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::runtime:list",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime-list": [{
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
