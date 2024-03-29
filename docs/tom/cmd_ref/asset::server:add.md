# DESCRIPTION

This command is used to create a new server within the specified
namespace. A server must have a type that is either `virtual` or
`physical`. Physical servers have no parent entity that provides them.
Virtual servers are provided by runtime environments and can be stacked
on top of them.

# SYNOPSIS

```
tom server add ${name} namespace ${space} type ${typ} [since ${since}] [until ${until}] [property ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
typ | string | type of the server | | no

# NOTES

The time specification for `since` can be either the special keyword
`always` or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to -4096-01-01T00:00:00Z, which is -infinity
for the system.

The time specification for `until` can be either the special keyword
`forever` or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to 293888-01-01T00:00:00Z, which is +infinity
for the system.

The `since` and `until` specification of the runtime are set as the
validity of the `name` property. After its creation, the `since`
timestamp can not be updated to earlier points in time.

If unspecified, the default value for `since` is the time of the command,
the default value for `until` is forever.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server add example-db01 namespace inventory type physical since always until forever property interface_list value '["lo"]' property if_lo_link_addr_list value '["00:00:00:ff:ff:ff"]'
```
