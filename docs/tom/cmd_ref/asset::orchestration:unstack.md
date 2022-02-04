# DESCRIPTION

This command can be used to remove a provider from an orchestration
environment.

# SYNOPSIS

```
tom orchestration unstack ${name} namespace ${space} unprovide ${providerID} [unprovide ...]
tom orchestration unstack ${tomID} unprovide ${providerID} [unprovide ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the orchestration | | no
space | string | name of the namespace | | no
tomID | string | tomID of the orchestration | | no
providerID | string | tomID of the entity for removal | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom orchestration unstack k8s-cluster.inventory.orchestration.tom unprovide k8s-cluster-01a_OLD.inventory.runtime.tom
```
