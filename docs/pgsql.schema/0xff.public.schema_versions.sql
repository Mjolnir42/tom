--
--
-- PUBLIC SCHEMA
CREATE TABLE IF NOT EXISTS  public.schema_versions (
    serial                        bigserial       PRIMARY KEY,
    schema                        varchar(16)     NOT NULL,
    version                       numeric(16,0)   NOT NULL,
    created_at                    timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    description                   text            NOT NULL
);

CREATE VIEW view.schema_version AS
SELECT  schema AS schema,
        MAX(version) AS version
FROM    public.schema_versions
GROUP   BY schema;

INSERT INTO public.schema_versions (
    schema,
    version,
    description )
VALUES
(   'meta',
    20191010001,
    'initial schema installation'
),
(   'ix',
    20191010001,
    'initial schema installation'
),
(   'filter',
    20191010001,
    'initial schema installation'
),
(   'yp',
    20191010001,
    'initial schema installation'
),
(   'asset',
    20191010001,
    'initial schema installation'
),
(   'view',
    20191010001,
    'initial schema installation'
),
(   'bulk',
    20191010001,
    'initial schema installation'
),
(   'filter',
    20191011001,
    'CIAA mapping update'
),
(   'ix',
    20191011001,
    'rename information_system_component to functional_component'
),
(   'view',
    20191011001,
    'rename views so they show up in \d'
),
(   'ix',
    20191014001,
    'rename logical_component_subgroup to deployment_group'
),
(   'filter',
    20191014001,
    'adopt deployment group renaming from schema ix'
),
(   'view',
    20191014001,
    'adopt deployment group renaming from schema ix'
),
(   'ix',
    20191015001,
    'use better names for attribute value assignment tables'
),
(   'meta',
    20191016001,
    'use better names for attribute value assignment tables'
),
(   'ix',
    20191017001,
    'adpopt long value table names'
),
(   'meta',
    20191017001,
    'adpopt long value table names'
),
(   'asset',
    20191017001,
    'adpopt long value table names'
),
(   'filter',
    20191017001,
    'adpopt long value table names'
),
(   'yp',
    20191017001,
    'adpopt long value table names'
),
(   'bulk',
    20191017001,
    'align table naming'
),
(
    'meta',
    20191022001,
    'fix constraint __mdq_temp_uniq'
),
(
    'asset',
    20191104001,
    'add server_parent table'
),
(
    'view',
    20191104001,
    'add dictionary schema/definition functions'
),
(
    'view',
    20191104002,
    'add resolveRuntimeTo.. functions'
),
(
    'filter',
    20191126001,
    'fully re-design filter schema'
),
(
    'asset',
    20200908001,
    'add tables for socket in schema asset'
),
(
    'ix',
    20200914001,
    'add endpoint tables'
),
(
    'filter',
    20200914001,
    'add endpoint as filter-able entity'
),
(
    'yp',
    20200915001,
    'rename schema tosm to yp'
),
(
    'yp',
    20200915002,
    'add service tables'
),
(
    'bulk',
    20200915001,
    'rename instance data table'
),
(
    'asset',
    20200915001,
    'remove orchestration environments as possible socket parent'
),
(
    'ix',
    20200915001,
    'rename technical system services to technical services'
),
(
    'yp',
    '20200917001',
    'cleanup of relationship tables',
),
(
    'asset',
    '20200917001',
    'cleanup of relationship tables',
),
(
    'ix',
    20200917001,
    'cleanup of relationship tables'
),
(
    'asset',
    20201015001,
    'add container entity support'
),
(
    'yp',
    20201016001,
    'switch from service linking 1:n to mapping n:m'
),
(
    'filter',
    20201016001,
    'add container as filter-able entity'
)
;
