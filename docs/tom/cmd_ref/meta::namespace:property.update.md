# DESCRIPTION

This command is used to update one or more properties of a namespace in a
single call. All specified properties are updated if the value changes.
All previously set properties not specified in the command are left at
their original value.

At least one property has to be specified.

The validity of the value can be specified using the optional since and
until keywords.

Not previously created attributes specified in the call are created as
standard attributes.

# SYNOPSIS

```
tom namespace property update ${name} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the namespace | | no
attr | string | name of the property attribute | | no
val | string | value of the property | | no
since | string | since when the value is valid for the property | now | yes
until | string | until when the value is valid for the property | forever | yes

# NOTES

The time specification for `since` can be either the special keyword
'always' or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to -4096-01-01T00:00:00Z.

The time specification for `until` can be either the special keyword
'forever' or a timestamp in RFC3339 format with millisecond precision.
The keyword translates to 293888-01-01T00:00:00Z.

An property value update can not move the since validity of the new
value before the since validity of old value.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom namespace property set inventory property foobar value since always until 2021-10-31T09:57:03.000+02:00
tom namespace property set inventory property foobar value test1 property barfoo value test2
```
