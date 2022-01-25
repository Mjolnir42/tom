# DESCRIPTION

This command is used to register which runtime provides a
server. It is only valid for servers of type `virtual`.
The providing entity must be a runtime environment.

The server can either be specified by name and namespace combination or
tomID. The runtime environment must be specified as tomID.
Either the DNS or URI form of the tomID is valid in either case.

# SYNOPSIS

```
tom server stack ${name} namespace ${space} provided-by ${providerID}
tom server stack ${tomID} provided-by ${providerID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
tomID | string | tomID of the server | | no
providerID | string | tomID of the runtime environment | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server stack example-db01 namespace inventory provided-by hypervisor01-kvm.inventory.runtime.tom
tom server stack example-db01.inventory.server.tom provided-by hypervisor01-kvm.inventory.runtime.tom
tom server stack tom://inventory/server/name=example-db01 provided-by tom://inventory/runtime/name=hypervisor01-kvm
```
