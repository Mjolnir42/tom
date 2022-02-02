# DESCRIPTION

This command is used to link two orchestration environments as referencing
the same real-world object.

The created link is perpetual and can not be split up, ie. it is
invalid that two orchestration environments were the exact same orchestration
environments at one point but later are not.

Linked orchestration environments can be from the same or different namespaces.

Orchestration environments have to be specified via tomID, which can be given in
either DNS or URI format.

# SYNOPSIS

```
tom orchestration link ${tomID} is-equal ${linkedID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
tomID | string | tomID of the first orchestration | | no
linkedID | string | tomID of the second orchestration | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom orchestration link kvm-block01a.inventory.orchestration.tom is-equal kvm-block01a.team~virt.orchestration.tom
tom orchestration link tom://inventory/orchestration/name=kvm-block01a is-equal tom://team~virt/orchestration/name=kvm-block01a
```
