# COMMAND

`asset::server:unstack`

# REQUEST PATH

```
DELETE /server/:tomID/parent
```

# REQUEST BODY

```
none
```

# RESPONSE

```
{ "command":   "asset::server:link",
  "error":     "",
  "requestID": "<uuid>",
  "status":    200,
  "server": [{
    "namespace": "<string>",
    "name": "<string>",
  }]
}
```
