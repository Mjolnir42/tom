# runtime DEFINITION

Runtime environments are the general purpose object for software
stacks on compute resources.

The `tomID` of a runtime environment ends in `.rte.tom` or `.runtime.tom`.

# SYNOPSIS OVERVIEW

```
tom runtime list [namespace ${space}]
tom runtime show ${name} namespace ${space}
tom runtime show ${tomID}
tom runtime add ${name} namespace ${space} type ${typ} [since ${since}] [until ${until}] [property ...]
tom runtime remove ${name} namespace ${space}
tom runtime remove ${tomID}
tom runtime property set ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom runtime property update ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom runtime property remove ${name} namespace ${space} property ${attr} [property ...]
tom runtime link ${tomID} is-equal ${linkedID}
tom runtime stack ${name} namespace ${space} provided-by ${providerID}
tom runtime stack ${tomID} provided-by ${providerID}
tom runtime resolve ${name} namespace ${space} level ${detail}
tom runtime resolve ${tomID} level ${detail}
tom runtime unstack ${name} namespace ${space}
tom runtime unstack ${tomID}
```

# sys PROPERTIES

The following are the properties a runtime environment should have in
namespace `sys`.
Perpetual properties can not be changed, while properties with validity
can be updated over the lifetime of the runtime environment.

Attribute | Unique? | Perpetual
 -------- | ------- | ---------
name | yes | no
bios_uuid | yes | no
type | no | yes
parent | no | no
owner | no | no
lifecycle_state | no | no
fqdn | no | no
interface_list | no | no
if_{}_link_addr_list | no | no
if_{}_ip_addr_list | no | no
listen_socket_list | no | no
os_release | no | no

For every item in `interface_list`, entries `if_{}_link_addr_list` and
`if_{}_ip_addr_list` should exist with the `{}` placeholder replaced
with the item.

# SEE ALSO

See `tom runtime help ${command}` for detailed help.
