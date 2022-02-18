# COMMAND

`asset::runtime:resolve`

# REQUEST PATH

```
GET /runtime/:tomID/resolve/:level
```

The `level` argument can be `server` to resolve to any server type or
`physical` to resolve down to physical servers.

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::runtime:resolve",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime-list": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
  },
  ...
  ]
}
```
