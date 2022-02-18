# COMMAND

`asset::server:resolve`

# REQUEST PATH

```
GET /server/:tomID/resolve/:level
```

The `level` argument can be `server` to resolve to any server type or
`physical` to resolve down to physical servers.


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
