# COMMAND

`asset::container:unstack`

# REQUEST PATH

```
DELETE /container/:tomID/parent
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::container:unstack",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
    "namespace": "<string>",
    "name":      "<string>",
  }]
}
```
