# DESCRIPTION

This command retrieves a list of configured runtime environments.
The optional namespace parameter restricts the resultset to the
specified namespace.

# SYNOPSIS

```
tom runtime list [namespace ${space}]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
space | string | name of the namespace | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom runtime list
tom runtime list namespace inventory
```
