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
[-] Parent
[x] Mapping     n:m Module
### Module
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Deployment
[-] Parent
[x] Mapping     n:m Artifact
### Artifact
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Instance
[-] Parent
[-] Mapping
### Data
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Shard
[x] Parent      n:1 Artifact
                n:1 Module
                n:1 Blueprint
[-] Mapping
### Service
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Realization 1:n Endpoint
[x] Parent      n:1 Blueprint
                n:1 Module
                n:1 Artifact
[-] Mapping

## Production Domain

### Technical Product
[x] Base
[x] SA
[x] QA
[ ] Linking
[-] Parent
[x] Mapping       n:m Deployment
### Deployment
[x] Base
[x] SA
[x] QA
[ ] Linking
[-] Parent
[x] Mapping       n:m Instance
### Instance
[x] Base
[x] SA
[x] QA
[ ] Linking
[-] Parent
[-] Mapping
### Shard
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Parent        n:1 Technical Product
                  n:1 Deployment
                  n:1 Instance
[-] Mapping
### Endpoint
[x] Base
[x] SA
[x] QA
[ ] Linking
[x] Parent        n:1 Technical Product
                  n:1 Deployment
                  n:1 Instance
[-] Mapping
### Netrange
[x] Base
[x] SA
[x] QA
[ ] Linking
[-] Parent
[x] Mapping       n:m Technical Product
                  n:m Deployment
                  n:m Instance

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
