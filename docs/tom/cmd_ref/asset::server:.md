# server DEFINITION

Servers are packages of compute resources. All servers must have a type
that is either `physical` or `virtual`.
For physical servers, the package corresponds to a physical server
chassis with hardware built into it. For virtual servers, the package
corresponds to a set of exclusive resources provided by a hypervisor
software.

# SYNOPSIS OVERVIEW

```
tom server add ${name} namespace ${space} type ${typ} [since ${since}] [until ${until}] [property ...]
tom server list [namespace ${space}]
tom server show ${name} namespace ${space}
tom server show ${tomID}
tom server property set ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom server property update ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom server property remove ${name} namespace ${space} property ${attr} [property ...]
tom server link ${tomID} is-equal ${linkedID}
tom server stack ${name} namespace ${space} provided-by ${providerID}
tom server stack ${tomID} provided-by ${providerID}
tom server unstack ${name} namespace ${space}
tom server unstack ${tomID}
```

# PROPERTIES

The following are the properties a server should have.
Perpetual properties can not be changed, while properties with validity
can be updated over the lifetime of the server.

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
