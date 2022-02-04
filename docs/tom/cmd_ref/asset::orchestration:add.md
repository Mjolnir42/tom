# DESCRIPTION

This command is used to create a new orchestration environment within the
specified namespace. Aside from its mandatory parameters, properties
of the orchestration can also be set in the same command.

If the optional since/until keywords are given, they should be placed
before the first property declaration, since the same keywords are
also used for properties.

Setting properties for unknown attributes creates those attributes in
the namespace as non-unique standard attributes.

# SYNOPSIS

```
tom orchestration add ${name} namespace ${space} type ${typ} [since ${since}] [until ${until}] [property ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the orchestration | | no
space | string | name of the namespace | | no
typ | string | type of the orchestration | | no
since | timestamp | since when this orchestration is valid | now | yes
until | timestamp | until when this orchestration is valid | forever | yes

# NOTES

The time specification for `since` can be either the special keyword
`always` or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to -4096-01-01T00:00:00Z, which is -infinity
for the system.

The time specification for `until` can be either the special keyword
`forever` or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to 293888-01-01T00:00:00Z, which is +infinity
for the system.

The `since` and `until` specification of the orchestration are set as the
validity of the `name` property. After its creation, the `since`
timestamp can not be updated to earlier points in time.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom orchestration add k8s-cluster namespace inventory type Kubernetes since always until forever
```
