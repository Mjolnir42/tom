# DESCRIPTION

This command is used to show full details about a server.
There are two variants of this command, referencing the server either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom server show ${name} namespace ${space}
tom server show ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
tomID | string | tomID of the server | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server show example-db01 namespace inventory
tom server show example-db01.inventory.server.tom
tom server show tom://inventory/server/name=example-db01
```
