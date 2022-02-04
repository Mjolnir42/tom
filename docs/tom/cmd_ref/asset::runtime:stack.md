# DESCRIPTION

This command is used to register which entity provides the current
runtime. The providing entity must be a runtime environment, a server or
an orchestration environment. Containers can not be providing entities.

The server can either be specified by name and namespace combination or
tomID. The runtime environment must be specified as tomID.
Either the DNS or URI form of the tomID is valid in either case.

# SYNOPSIS

```
tom runtime stack ${name} namespace ${space} provided-by ${providerID}
tom runtime stack ${tomID} provided-by ${providerID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the runtime | | no
space | string | name of the namespace | | no
tomID | string | tomID of the server | | no
providerID | string | tomID of the providing entity | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom runtime stack example-db01 namespace inventory provided-by example-db01.inventory.server.tom
tom runtime stack hypervisor01-kvm.inventory.runtime.tom provided-by hypervisor01-virt.inventory.orchestration.tom
tom runtime stack hypervisor01-chroot.inventory.runtime.tom provided-by hypervisor01.inventory.runtime.tom
```
