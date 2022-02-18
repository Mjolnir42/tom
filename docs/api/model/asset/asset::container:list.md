# COMMAND

`asset::container:list`

# REQUEST PATH

```
GET /container/
GET /container/?namespace=<string>
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::container:list",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container-list": [{
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
