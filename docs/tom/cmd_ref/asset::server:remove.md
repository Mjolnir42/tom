# DESCRIPTION

The remove command can be used to delete a server. Removal is a soft delete,
that ends the validity of all records describing the server.

Upon removal, the server is unstacked from its parent. The server also
unstacks all its children from itself.

# SYNOPSIS

```
tom server remove ${name} namespace ${space}
tom server remove ${tomID}
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the server | | no
space | string | name of the namespace | | no
tomID | string | TomID of the server | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom server remove example-db01 namespace inventory
```
