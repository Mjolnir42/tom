# DESCRIPTION

This command is used to link two servers as referencing
the same real-world object.

The created link is perpetual and can not be split up, ie. it is
invalid that two servers were the exact same server at one point
but later are not.

Linked servers can be from the same or different namespaces.

Servers have to be specified via tomID, which can be given in
either DNS or URI format.

# SYNOPSIS

```
tom server link ${tomID} is-equal ${linkedID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
tomID | string | tomID of the first server | | no
linkedID | string | tomID of the second server | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server link example-db01.inventory.server.tom is-equal example-db01.team~dba.server.tom
tom server link tom://inventory/server/name=example-db01 is-equal tom://team~dba/server/name=example-db01
```
