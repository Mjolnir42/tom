# namespace DEFINITION

Namespaces are the scoping structure of TOM. Every data entry is associated
with a namespace.

Namespaces are also the boundaries of data models, as each namespace has its
own possible attributes to construct it.

By convention, every namespace must have a unique `name` and a standard `type`
attribute.

Namespace names may only consist of the following characters: `a-zA-z0-9_-`.
In addition, a namespace name may have one of the following prefixes:
`team~`, `tool~`.

# SYNOPSIS OVERVIEW

```
tom namespace add ${name} type ${type} [lookup-uri ${uri}] [lookup-key ${key}] [entities ${ntt}] [std-attr ${std}] [uniq-attr ${uniq}]
tom namespace list
tom namespace show ${name}
tom namespace attribute add ${name} [std-attr ${std}] [uniq-attr ${uniq}]
tom namespace attribute remove ${name} [std-attr ${std}] [uniq-attr ${uniq}]
tom namespace remove ${name}
```

See `tom namespace help ${command}` for detailed help.
