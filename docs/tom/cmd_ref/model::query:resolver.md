# DESCRIPTION

This command is used to either show or list information and details based on the presented single argument. That single argument can either be a full tomID, or one of the supported wildcard formats.

The tomID parameter can be given in either DNS or URI format.

The wildcard parameter only supports the DNS format style.

# SYNOPSIS

```
tom query ${tomID}
tom query ${wildcard}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
tomID | string | tomID of the entity | | no
wildcard | string | wildcard ID of the entities | | no

# NOTES

For wildcard requests to work, shell globbing has to be disabled for the `tom` command.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom query example-db01.inv.srv.tom
tom query *.inv.srv.tom
tom query *.srv.tom
```
