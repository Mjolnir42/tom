# COMMAND

`asset::runtime:unstack`

# REQUEST PATH

```
DELETE /runtime/:tomID/parent
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::runtime:unstack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "runtime": [{
    "namespace": "<string>",
    "name":      "<string>",
  }]
}
```
