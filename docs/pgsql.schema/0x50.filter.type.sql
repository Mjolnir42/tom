---
---
--- FILTER SCHEMA
CREATE TYPE flt_card AS ENUM(
    'one',
    'many'
);
CREATE TYPE flt_aggr AS ENUM(
    'min',
    'max',
    'first',
    'last'
);
CREATE TYPE flt_ntt AS ENUM(
    'top_level_service',
    'product',
    'information_system',
    'functional_component',
    'deployment_group',
    'runtime_environment',
    'orchestration_environment',
    'server',
    'endpoint',
    'container'
);
