# DESCRIPTION

This command retrieves a list of configured server. The optional
namespace parameter restricts the resultset to the specified namespace.

# SYNOPSIS

```
tom server list [namespace ${space}]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
namespace | string | namespace to list | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server list
tom server list namespace inventory
```
