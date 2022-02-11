# DESCRIPTION

This command can be used to remove the provider from a runtime
environment.

It is not required to know the current provider of the runtime.

# SYNOPSIS

```
tom runtime unstack ${name} namespace ${space}
tom runtime unstack ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the runtime | | no
space | string | name of the namespace | | no
tomID | string | tomID of the runtime | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom runtime unstack example-db01 namespace inventory
```
