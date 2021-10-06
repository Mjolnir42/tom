# DESCRIPTION

This command is used to remove a runtime environment. This is done by
ending the runtime environment's validity timestamp, ie. the runtime environment is
not fully deleted and will continue to show up when the past
is queried.

There are two variants of this command, referencing the runtime environment either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom runtime remove ${name} namespace ${space}
tom runtime remove ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the runtime | | no
space | string | name of the namespace | | no
tomID | string | tomID of the runtime | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom runtime remove example-db01 namespace inventory
tom runtime remove example-db01.inventory.runtime.tom
tom runtime remove tom://inventory/runtime/name=example-db01
```
