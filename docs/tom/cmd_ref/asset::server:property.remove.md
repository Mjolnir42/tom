# DESCRIPTION

This command is used to remove one or more properties of a server
in a single call.
It does not delete the current value from the system, but updates
its validity to end.

More than one property can be specified in a single command. The is not
required to currently have a value, but the attribute must exist within
the namespace.

At least one property has to be specified.

# SYNOPSIS

```
tom server property remove ${name} namespace ${space} property ${attr} [property ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
attr | string | name of the property attribute | | no

# NOTES

The following attributes are never removed by this call:

Attribute | Reason
 -------- | ------
name | Invalidating the name is functionally deleting the server.
type | The type is perpetual and can not be changed.
parent | Inventory stack information.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server property remove example-db01 namespace inventory property lifecycle_state
```
