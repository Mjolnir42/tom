# COMMAND

`asset::server:resolve`

# REQUEST PATH

```
GET /server/:tomID/resolve/:level
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::server:resolve",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server-list": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
  },
  ...
  ]
}
```
