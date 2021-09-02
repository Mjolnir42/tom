# namespace DEFINITION

Namespaces are the scoping structure of TOM. Every data entry is associated
with a namespace.

# SYNOPSIS OVERVIEW

```
tom namespace add ${name} type ${type} [lookup-uri ${uri}] [lookup-key ${key}] [entities ${ntt}] [std-attr ${std}] [uniq-attr ${uniq}]
tom namespace list
tom namespace show ${name}
tom namespace attribute add ${name} [std-attr ${std}] [uniq-attr ${uniq}]
tom namespace attribute remove ${name} [std-attr ${std}] [uniq-attr ${uniq}]
```

See `tom namespace help ${command}` for detailed help.
