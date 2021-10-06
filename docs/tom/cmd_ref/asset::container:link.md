# DESCRIPTION

This command is used to link two containers as referencing the
same real-world object.

The created link is perpetual and can not be split up, ie. it is
invalid that two containers were the exact same container at one
point but later are not.

Linked containers can be from the same or different namespaces.

Containers have to be specified via tomID, which can be given in
either DNS or URI format.

# SYNOPSIS

```
tom container link ${tomID} is-equal ${linkedID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
tomID | string | tomID of the first container | | no
linkedID | string | tomID of the second container | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom container link example-db01.inventory.container.tom is-equal example-db01.team~dba.container.tom
tom container link tom://inventory/container/name=example-db01 is-equal tom://team~dba/container/name=example-db01
```
