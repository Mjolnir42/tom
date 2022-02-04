# DESCRIPTION

This command is used to register which entities provides the current
orchestration. The providing entity must be a runtime environment.

The orchestration can either be specified by name and namespace combination or
tomID. The runtime environment must be specified as tomID.
Either the DNS or URI form of the tomID is valid in either case.

Since orchestration environments have multiple providers, this call only
adds providers to the already registered ones. Multiple `provided-by`
clauses can be added in one command.
If providers are to be deregistered, the `replacing` clause can be used
to indicate which providers to remove. Multiple `replacing` clauses can
be used in one command as well.

# SYNOPSIS

```
tom orchestration stack ${name} namespace ${space} provided-by ${providerID} [provided-by ...]
tom orchestration stack ${tomID} provided-by ${providerID} [provided-by ...] [replacing ${oldID}] [replacing ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the orchestration | | no
space | string | name of the namespace | | no
tomID | string | tomID of the orchestration | | no
providerID | string | tomID of the providing runtime | | yes
oldID | string | tomID of the runtime to be removed | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom orchestration stack k8s-cluster namespace inventory provided-by
k8s-cluster-01a.inventory.runtime.tom
tom orchestration stack k8s-cluster.inventory.orchestration.tom provided-by k8s-cluster-01a.inventory.runtime.tom
tom orchestration stack k8s-cluster.inventory.orchestration.tom provided-by k8s-cluster-01a.inventory.runtime.tom replacing k8s-cluster-01a_OLD.inventory.runtime.tom
```
