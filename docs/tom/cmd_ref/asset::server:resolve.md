# DESCRIPTION

This command can be used to resolve a server down to its
providing server. The requested detail level can either be `server`, in
which case the next parent(s) of type server are returned (either virtual or
physical). If the requested detail level is `physical`, then the stack
is followed down until the stacking terminates at physical servers.

The result set can contain multiple servers if an orchestration
environment is used within the stack.

The result for an physical server is always itself. The result for a virtual
server is itself or its physical server, depending on requested detail level.

# SYNOPSIS

```
tom server resolve ${name} namespace ${space} level ${detail}
tom server resolve ${tomID} level ${detail}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
tomID | string | tomID of the server | | no
detail | string | requested detail level | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server resolve example-db01 namespace inventory level server
tom server resolve example-db01.inventory.server.tom level physical
```
