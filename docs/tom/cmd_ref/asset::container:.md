# container DEFINITION

Container are the specially typed object to differentiate them
from full general purpose software stacks on compute resources.

# SYNOPSIS OVERVIEW

```
tom container list [namespace ${space}]
tom container show ${name} namespace ${space}
tom container show ${tomID}
tom container add ${name} namespace ${space} type ${typ} [since ${since}] [until ${until}] [property ...]
tom container remove ${name} namespace ${space}
tom container remove ${tomID}
tom container property set ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom container property update ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom container property remove ${name} namespace ${space} property ${attr} [property ...]
tom container link ${tomID} is-equal ${linkedID}
tom container stack ${name} namespace ${space} provided-by ${providerID}
tom container stack ${tomID} provided-by ${providerID}
tom container unstack ${name} namespace ${space}
tom container unstack ${tomID}
tom container resolve ${name} namespace ${space} level ${detail}
tom container resolve ${tomID} level ${detail}
```

# PROPERTIES

The following are the properties a container should have.
Perpetual properties can not be changed, while properties with validity
can be updated over the lifetime of the container.

Attribute | Unique? | Perpetual
 -------- | ------- | ---------
name | yes | no
type | no | yes
lifecycle_state | no | no
base_image | no | no

# SEE ALSO

See `tom container help ${command}` for detailed help.
