# Authoritative Data

## Asset Domain

### Server
[ ] Done.
### Runtime Environment
[ ] Done.
### Orchestration Environment
[ ] Done.
### Container
[ ] Done.
### Socket
[ ] ToDo


## Abstract Service Domain

### Blueprint
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n TechnicalProduct
[x] Mapping     n:m Module
### Module
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Deployment
[x] Mapping     n:m Artifact
### Artifact
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Instance
### Data
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Shard
[x] Parent      n:1 Artifact
                n:1 Module
                n:1 Blueprint
### Service
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Endpoint
[x] Parent      n:1 Blueprint
                n:1 Module
                n:1 Artifact

## Production Domain

### Technical Product (Blueprint)
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Mapping       n:m Deployment
    [ ] :add
    [ ] :list
    [ ] :remove
    [ ] :show
    [ ] :contains Deployment
### Deployment  (Module)
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Mapping       n:m Instance
    [ ] :add
    [ ] :list
    [ ] :remove
    [ ] :show
    [ ] :partof Deployment
### Instance (Artifact)
[x] Base
[x] SA
[x] QA
[ ] Linking
    [ ] :add
    [ ] :list
    [ ] :remove
    [ ] :show
    [ ] :connects-to  Endpoint
### Shard (Data)
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Parent        n:1 Technical Product
                  n:1 Deployment
                  n:1 Instance
### Endpoint (Service)
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Parent        n:1 Technical Product
                  n:1 Deployment
                  n:1 Instance
    [ ] :add
    [ ] :list
    [ ] :remove
    [ ] :show
    [ ] :provided-by  asset::*
### Netrange
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Mapping       n:m Technical Product
                  n:m Deployment
                  n:m Instance
    [ ] :add
    [ ] :list
    [ ] :remove
    [ ] :show
    [ ] :uses  Endpoint

## Reporting Domain

### Consumer Product
[x] Base
[x] SA
[x] QA
[ ] Linking
[ ] Parent
[ ] Mapping
### Top Level Service
[x] Base
[x] SA
[x] QA
[ ] Linking
[ ] Parent
[ ] Mapping

# Referential Data

## Related Links

### Corporate Domain          [DOMAIN- ]
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Parent
[ ] Mapping
### Domain                    [DOMAIN- ]
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Parent
[ ] Mapping
### Information System        [IS-     ]
[x] Base
[x] SA
[x] QA
[ ] Linking
[ ] Parent
[ ] Mapping
### Service                   [SER-    ]
[x] Base
[x] SA
[x] QA
[ ] Linking
[ ] Parent
[ ] Mapping
### Software Asset            [YP-     ]
[x] Base
[x] SA
[x] QA
[ ] Linking
[ ] Parent
[ ] Mapping
### Technology Reference Card [TRC-    ]
[x] Base
[x] SA
[x] QA
[ ] Linking
[ ] Parent
[ ] Mapping
