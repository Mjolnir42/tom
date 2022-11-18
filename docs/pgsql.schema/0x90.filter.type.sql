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
    'blueprint',
    'module',
    'artifact',
    'data',
    'service',
    'technical_product',
    'deployment',
    'instance',
    'shard',
    'endpoint',
    'netrange',
    'consumer_product',
    'top_level_service',
    'server',
    'runtime_environment',
    'orchestration_environment',
    'container'
);
