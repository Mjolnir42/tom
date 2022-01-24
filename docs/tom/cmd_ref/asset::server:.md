# server DEFINITION

Servers are packages of compute resources. All servers must have a type
that is either `physical` or `virtual`.
For physical servers, the package corresponds to a physical server
chassis with hardware built into it. For virtual servers, the package
corresponds to a set of exclusive resources provided by a hypervisor
software.

# SYNOPSIS OVERVIEW

```
tom server link ${tomID} is-equal ${linkedID}
```

# PROPERTIES

The following are the properties a server should have.
Perpetual properties can not be changed, while properties with validity
can be updated over the lifetime of the runtime environment.

Attribute | Unique? | Perpetual
 -------- | ------- | ---------
name | yes | no
type | no | yes
interface_list | no | no
if_{}_link_addr_list | no | no

For every item in `interface_list`, an entry in `if_{}_link_addr_list`
should exist with the `{}` placeholder replaced with the item.

# SEE ALSO

See `tom server help ${command}` for detailed help.
