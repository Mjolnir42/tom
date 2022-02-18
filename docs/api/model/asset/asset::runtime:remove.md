# COMMAND

`asset::runtime:remove`

# REQUEST PATH

```
DELETE /runtime/:tomID
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::runtime:remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
    "namespace": "<string>",
    "name":      "<string>"
  }]
}
```
