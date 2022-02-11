# orchestration DEFINITION

Orchestration environments are the multi-parent objects that can be used
to model clusters with their own resource orchestration. If there is a
scheduling layer that autonomously moves resources within a cluster, and
it is not desirable to replicate every individual movement into `Tom`,
then an orchestration environment is of use. Everything above it is
provided-by everything below it.

# SYNOPSIS OVERVIEW

```
tom orchestration list [namespace ${space}]
tom orchestration show ${name} namespace ${space}
tom orchestration show ${tomID}
tom orchestration add ${name} namespace ${space} type ${typ} [since ${since}] [until ${until}] [property ...]
tom orchestration remove ${name} namespace ${space}
tom orchestration remove ${tomID}
tom orchestration property update ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom orchestration property set ${name} namespace ${space} property ${attr} value ${val} [since ${since}] [until ${until}] [property ...]
tom orchestration property remove ${name} namespace ${space} property ${attr} [property ...]
tom orchestration link ${tomID} is-equal ${linkedID}
tom orchestration stack ${name} namespace ${space} provided-by ${providerID} [provided-by ...]
tom orchestration stack ${tomID} provided-by ${providerID} [provided-by ...] [replacing ${oldID}] [replacing ...]
tom orchestration unstack ${name} namespace ${space} unprovide ${providerID} [unprovide ...]
tom orchestration unstack ${tomID} unprovide ${providerID} [unprovide ...]
tom orchestration resolve ${name} namespace ${space} level ${detail}
tom orchestration resolve ${tomID} level ${detail}
```

# PROPERTIES

The following are the properties an orchestration environment should have.
Perpetual properties can not be changed, while properties with validity
can be updated over the lifetime of the orchestration environment.

Attribute | Unique? | Perpetual
 -------- | ------- | ---------
name | yes | no
type | no | yes

# SEE ALSO

See `tom orchestration help ${command}` for detailed help.
