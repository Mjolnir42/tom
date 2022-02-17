# COMMAND

`asset::server:remove`

# REQUEST PATH

```
DELETE /server/:tomID
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::server:remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
    "namespace": "<string>",
    "name":      "<string>"
  }]
}
```
