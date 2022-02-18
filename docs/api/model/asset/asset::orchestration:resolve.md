# COMMAND

`asset::orchestration:resolve`

# REQUEST PATH

```
GET /orchestration/:tomID/resolve/:level
```

The `level` argument can be `server` to resolve to any server type or
`physical` to resolve down to physical servers.

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::orchestration:resolve",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration-list": [{
    "namespace": "<string>",
    "name":      "<string>",
    "type":      "<string>",
  },
  ...
  ]
}
```
