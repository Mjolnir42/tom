# COMMAND

`asset::orchestration:remove`

# REQUEST PATH

```
DELETE /orchestration/:tomID
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::orchestration:remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name":      "<string>"
  }]
}
```
