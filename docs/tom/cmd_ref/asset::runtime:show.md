# DESCRIPTION

This command is used to show full details about a runtime environment.
There are two variants of this command, referencing the runtime either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom runtime show ${name} namespace ${space}
tom runtime show ${tomID}
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
tom runtime show example-db01 namespace inventory
tom runtime show example-db01.inventory.runtime.tom
tom runtime show tom://inventory/runtime/name=example-db01
```
