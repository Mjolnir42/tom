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
{ "command":   "asset::orchestration:link",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "orchestration": [{
    "namespace": "<string>",
    "name": "<string>",
  }]
}
```
