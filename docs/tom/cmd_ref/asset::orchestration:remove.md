# DESCRIPTION

This command is used to remove a orchestration environment. This is done by
ending the orchestration environment's validity timestamp, ie. the orchestration environment is
not fully deleted and will continue to show up when the past
is queried.

There are two variants of this command, referencing the orchestration environment either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom orchestration remove ${name} namespace ${space}
tom orchestration remove ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the orchestration | | no
space | string | name of the namespace | | no
tomID | string | tomID of the orchestration | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom orchestration remove k8s-cluster namespace inventory
tom orchestration remove k8s-cluster.inventory.orchestration.tom
tom orchestration remove tom://inventory/orchestration/name=k8s-cluster
```
