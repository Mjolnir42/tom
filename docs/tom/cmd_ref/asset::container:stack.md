# DESCRIPTION

This command is used to register which runtime provides a container.
The providing entity must be a runtime environment.

The container can either be specified by name and namespace combination or
tomID. The runtime environment must be specified as tomID.
Either the DNS or URI form of the tomID is valid in either case.

# SYNOPSIS

```
tom container stack ${name} namespace ${space} provided-by ${providerID}
tom container stack ${tomID} provided-by ${providerID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the container | | no
space | string | name of the namespace | | no
tomID | string | tomID of the container | | no
providerID | string | tomID of the runtime environment | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom container stack foocontainer namespace inventory provided-by k8s.inventory.runtime.tom
```
