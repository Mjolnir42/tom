# DESCRIPTION

This command is used to remove one or more properties of a namespace in a
single call. It does not delete the current value from the system, but
updates its validity to end.

More than one property can be specified in a single command. The is not
required to currently have a value, but the attribute must exist within
the namespace.

# SYNOPSIS

```
tom namespace property remove ${namespace} property ${attribute} [property ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
namespace | string | name of the namespace | | no
attribute | string | name of the property attribute | | no

# NOTES

All properties removed in the same call are invalidated using the same
timestamp.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom namespace property remove inventory property foobar
tom namespace property remove inventory property foobar property barfoo
```
