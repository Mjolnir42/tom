# DESCRIPTION

This command is used to show full details about a orchestration environment.
There are two variants of this command, referencing the orchestration either
by name and namespace combination, or by tomID.
The tomID parameter can be given in either DNS or URI format.

# SYNOPSIS

```
tom orchestration show ${name} namespace ${space}
tom orchestration show ${tomID}
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
tom orchestration show k8s-cluster namespace inventory
tom orchestration show k8s-cluster.inventory.orchestration.tom
tom orchestration show tom://inventory/orchestration/name=k8s-cluster
```
