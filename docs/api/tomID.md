# tomID

The `tomID` is the unique identifier of an object known to Tom. It is human
readable, an consistent only of URI safe characters.

## uri format

The general URI format of the `tomID` is:

```
tom://${namespace}/${entity}/name=${name}
```

For namespaces themselves, this results in:

```
tom:///${entity}/name=${name}
```


## dns format

The general DNS format of the `tomID` is:

```
${name}.${namespace}.${entity}.tom
${name}.${namespace}.${entity}.tom.
```

Keeping with tradition, the DNS format may end with a `.` character that can
be omitted.

The wildcard DNS format versions are understood by the cli utility `tom`,
but not used on the api level:

```
*.${namespace}.${entity}.tom
*.${entity}.tom
```

## entity specifiers

The following entity specifiers are defined. The short-version identifiers
are only valid for DNS formatted `tomID` strings.

1. Namespace: `namespace`
2. Server: `server`, `srv`
3. Runtime Environment: `runtime`, `rte`
4. Orchestration Environment: `orchestration`, `ore`
5. Container: `container`, `cnr`
6. Socket: `socket`, `sok`

# permitted characters

## permitted characters for namespaces

```
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~_-
```

The `~` character is only allowed to be used once, and only as part of the
following name prefixes:

1. `tool~`
2. `team~`

## permitted characters for object names

```
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~._-
```

## permitted characters for attribute names

```
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~.:_-
```

# examples

## namespace: inventory

```
inventory.namespace.tom
tom:///namespace/name=inventory
```

## server: git

```
git.inventory.server.tom
git.inventory.srv.tom
tom://inventory/server/name=git
```

## runtime: ftp01

```
ftp01.inventory.runtime.tom
ftp01.inventory.rte.tom
tom://inventory/runtime/name=ftp01
```

## orchestration: vz01-ha.example.com

```
vz01-ha.example.com.inventory.orchestration.tom.
vz01-ha.example.com.inventory.ore.tom
tom://inventory/orchestration/name=vz01-ha.example.com
```

## container: pgSQL13

```
pgSQL13.inventory.container.tom
pgSQL13.inventory.cnr.tom
tom://inventory/container/name=pgSQL13
```

## identity library: engineroom

```
engineroom.library.tom
tom:///library/name=engineroom
```

## machine: a7b13e3d8424a25a915f6fe9cfdd2b6a

```
a7b13e3d8424a25a915f6fe9cfdd2b6a.engineroom.machine.tom
tom://engineroom/machine/uid=a7b13e3d8424a25a915f6fe9cfdd2b6a
```
