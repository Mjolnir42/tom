# DESCRIPTION

This command is used to set specific properties of a runtime environment
in a single call. All specified properties are updated if the value changes.
See the `NOTES` section for properties that can not be updated.

All configured properties not specified in the command are left as is.

At least one property has to be specified.

The validity of the value can be specified using the optional since and
until keywords. Updating property values can not move the validity
further into the past.

Updates commands remove properties that have been set, but have not been
valid yet (ie. their valid since is still in the future) and would become
valid after the current command.

Specified, but not previously created attributes specified in the call
are transparently created as standard attributes.

# SYNOPSIS

```
tom runtime property update ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the runtime | | no
space | string | name of the namespace | | no
attr | string | name of the property attribute | | no
val | string | value of the property | | no
since | timestamp | since when this runtime is valid | now | yes
until | timestamp | until when this runtime is valid | forever | yes

# NOTES

The time specification for `since` can be either the special keyword
`always` or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to -4096-01-01T00:00:00Z, which is -infinity
for the system.
The keyword `now` as well as the unset default are the current time.

The time specification for `until` can be either the special keyword
`forever` or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to 293888-01-01T00:00:00Z, which is +infinity
for the system.
The default value when unset is `forever`. The keyword `now` is also
supported, but not very useful.

The `since` and `until` specification of the runtime are set as the
validity of the `name` property. After its creation, the `since`
timestamp can not be updated to earlier points in time, further into
the past.

The following attributes can not be updated:

Attribute | Reason
 -------- | ------
type | The type is perpetual and can not be changed.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom runtime property update example-db01 namespace inventory property lifecycle_state value end-of-life
```
