# DESCRIPTION

This command is used to remove a container. This is done by
ending the container's validity timestamp, ie. the container is
not fully deleted and will continue to show up when the past
is queried.

There are two variants of this command, referencing the container either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom container remove ${name} namespace ${space}
tom container remove ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the container | | no
space | string | name of the namespace | | no
tomID | string | tomID of the container | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom container remove example-db01 namespace inventory
tom container remove example-db01.inventory.container.tom
tom container remove tom://inventory/container/name=example-db01
```
