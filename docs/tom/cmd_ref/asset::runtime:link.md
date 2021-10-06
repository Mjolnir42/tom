# DESCRIPTION

This command is used to link two runtime environments as referencing
the same real-world object.

The created link is perpetual and can not be split up, ie. it is
invalid that two containers were the exact same container at one
point but later are not.

Linked runtime environments can be from the same or different namespaces.

Runtime environments have to be specified via tomID, which can be given in
either DNS or URI format.

# SYNOPSIS

```
tom runtime link ${tomID} is-equal ${linkedID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
tomID | string | tomID of the first runtime | | no
linkedID | string | tomID of the second runtime | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom runtime link example-db01.inventory.runtime.tom is-equal example-db01.team~dba.runtime.tom
tom runtime link tom://inventory/runtime/name=example-db01 is-equal tom://team~dba/runtime/name=example-db01
```
