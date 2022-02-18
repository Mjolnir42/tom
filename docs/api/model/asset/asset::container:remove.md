# COMMAND

`asset::container:remove`

# REQUEST PATH

```
DELETE /container/:tomID
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::container:remove",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "container": [{
    "namespace": "<string>",
    "name":      "<string>"
  }]
}
```
