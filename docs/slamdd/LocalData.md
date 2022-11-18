# Local data sync

```
/datadir  /sync   /namespace  /entity   /name
                                         ├── .since     now
                                         ├── .until     +1d
                                         ├── .style     update|set
                                         ├── .parent    tomID
                                         ├── .link      tomID....
                                         ├── attr       value
                                         ├── attr       value
                                         ├── attr       value
                                         └── attr       value

          /view   /namespace  /entity
                               ├── attr                 data/json
                               ├── attr                 data/json
                               └── name                 data/json
```
