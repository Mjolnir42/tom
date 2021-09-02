# DESCRIPTION

This command is used to remove an attribute from a namespace.
This is a *DESTRUCTIVE* debugging command that deletes all past, present and future
data associated with the attribute.

The intended usecase is to fix typos during namespace creation, nothing
more.

Multiple attributes can be purged at the same time.

# SYNOPSIS

```
tom namespace attribute remove ${name} [std-attr ${std}] [uniq-attr ${uniq}]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the namespace | | no
std | string | standard attribute of the namespace | | yes
uniq | string | unique attribute of the namespace | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
```
