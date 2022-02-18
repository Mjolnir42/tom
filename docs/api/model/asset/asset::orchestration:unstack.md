# COMMAND

`asset::orchestration:unstack`

# REQUEST PATH

```
DELETE /orchestration/:tomID/parent
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::orchestration:unstack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name":      "<string>",
  }]
}
```
