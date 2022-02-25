# DESCRIPTION

This command can be used to resolve a container down to its
providing server. The requested detail level can either be `server`, in
which case the next parent(s) of type server are returned (either virtual or
physical). If the requested detail level is `physical`, then the stack
is followed down until the stacking terminates at physical servers.

The result set can contain multiple servers if an orchestration
environment is used within the stack.

# SYNOPSIS

```
tom container resolve ${name} namespace ${space} level ${detail}
tom container resolve ${tomID} level ${detail}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the container | | no
space | string | name of the namespace | | no
tomID | string | tomID of the container | | no
detail | string | requested detail level | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# NOTES

Instead of `server`, the alias `next` is also recognized. Instead of
`physical`, the alias `full` is also recognized.

# EXAMPLES

```
tom container resolve dockerimg01 namespace inventory level server
tom container resolve dockerimg01.inventory.container.tom level physical
```
