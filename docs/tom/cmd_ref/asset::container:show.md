# DESCRIPTION

This command is used to show full details about a container.
There are two variants of this command, referencing the container either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom container show ${name} namespace ${space}
tom container show ${tomID}
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
tom container show example-db01 namespace inventory
tom container show example-db01.inventory.container.tom
tom container show tom://inventory/container/name=example-db01
```
