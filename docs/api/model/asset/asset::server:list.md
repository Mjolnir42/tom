# COMMAND

`asset::server:list`

# REQUEST PATH

```
GET /server/
GET /server/?namespace=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::server:list",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server-list": [{
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
