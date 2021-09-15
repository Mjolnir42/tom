# DESCRIPTION

This command is used to add a new attribute to a namespace. Multiple
attributes of same or different type can be added at the same time, but at
least one attribute must be specified.

# SYNOPSIS

```
tom namespace attribute add ${name} [std-attr ${std}] [uniq-attr ${uniq}]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the namespace | | no
std | string | Standard attribute for the namespace | | yes
uniq | string | Unique attribute for the namespace | | yes

# NOTES

If the name of the attribute ends in the suffix `_list`, then
its value must be decodeable by GoLang's encoding/json package
as array of strings.

If the name of the attribute ends in the suffix `_json`, then
its value must be decodeable by GoLang's encoding/json package
as a generic, arbitrary JSON.

If the name of the attribute ends in the suffix `_xml`, then
its value must be decodeable by GoLang's encoding/xml package.
as a generic, arbitrary XML.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
```
