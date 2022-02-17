# DESCRIPTION

This command is used to add a new namespace to TOM. Namespaces themselves are
identified by a `name`, and have a `type` attribute that must be either
`referential` or `authoritative`.

Referential namespaces only hold index information into data stored in a
different system. They must define a lookup key, which is the name of an
unique attribute in this namespace. The value of this unique attribute is
thereby declared as the unique identifier of the object in the other
system. A lookup URI must also be defined, which must contain the placeholder
string `{{LOOKUP}}`. If the value of the attribute referenced by
lookup-key is used to replace the placeholder in the lookup URI, then a
valid URI of the other system that handles a GET request and returns the
referenced object.
The referencing data might be augmented with additional relevant data points
required for correlation.

Authoritative namespaces hold data directly stored in TOM. Every object in a
namespace can have a list of associated key value pairs attached to it. The
keys of these key/value pairs are the attributes defined for the namespace.
Every object can have every attribute only once at the same time. For
standard attributes, multiple objects of the same type within the same
namespace can have the same value at the same time. For unique attributes, only
one object per type and namespace can have a specific value at any one time.

All namespaces must define a unique attribute `name` for use within them, as
all entity objects are referenced using this `name` attribute. All
namespaces must also define a standard attribute `type` for use within them.
These two attributes are always created for every namespace, even when not
specified.

# SYNOPSIS

```
tom namespace add ${name} type ${type} [lookup-uri ${uri}] [lookup-key ${key}] [entities ${ntt}] [std-attr ${std}] [uniq-attr ${uniq}]
```

# ARGUMENT TYPES

Argument | Type | Description | Default Value | Optional
 ------- | ---- | ----------- | ------------- | --------
name | string | name of the namespace | | no
type | string | type of the namespace | | no
uri | string | URI to use for lookups | | yes
key | string | Unique attribute used for lookups | | yes
ntt | list | UNSUPPORTED | | yes
std | string | Standard attribute for the namespace | | yes
uniq | string | Unique attribute for the namespace | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes

# EXAMPLES

```
tom namespace add inventory type authoritative uniq-attr name uniq-attr serial std-attr lifecycle-state
```
