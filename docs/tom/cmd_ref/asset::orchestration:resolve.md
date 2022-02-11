# DESCRIPTION

This command can be used to resolve an orchestration environment down to its
providing server. The requested detail level can either be `server`, in
which case the next parent(s) of type server are returned (either virtual or
physical). If the requested detail level is `physical`, then the stack
is followed down until the stacking terminates at physical servers.

The result set can contain multiple servers if an orchestration
environment is used within the stack.

# SYNOPSIS

```
tom orchestration resolve ${name} namespace ${space} level ${detail}
tom orchestration resolve ${tomID} level ${detail}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the orchestration | | no
space | string | name of the namespace | | no
tomID | string | tomID of the orchestration | | no
detail | string | requested detail level | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom orchestration resolve k8s namespace inventory level server
tom orchestration resolve k8s.inventory.orchestration.tom level physical
```
