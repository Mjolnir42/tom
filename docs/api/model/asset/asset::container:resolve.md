# COMMAND

`asset::container:resolve`

# REQUEST PATH

```
GET /container/:tomID/resolve/:level
```

The `level` argument can be `server` to resolve to any server type or
`physical` to resolve down to physical servers.

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::container:resolve",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container-list": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
  },
  ...
  ]
}
```
