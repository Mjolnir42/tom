# DESCRIPTION

This command can be used to remove the provider from a container.

It is not required to know the current provider.

# SYNOPSIS

```
tom container unstack ${name} namespace ${space}
tom container unstack ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
tomID | string | tomID of the server | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom container unstack example-db01 namespace inventory
```
