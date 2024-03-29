--
--
-- PUBLIC SCHEMA
CREATE TABLE IF NOT EXISTS  public.schema_versions (
    serial                        bigserial       PRIMARY KEY,
    schema                        varchar(16)     NOT NULL,
    version                       numeric(16,0)   NOT NULL,
    createdAt                     timestamptz(3)  NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
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
(   'meta',
    20191022001,
    'fix constraint __mdq_temp_uniq'
),
(   'asset',
    20191104001,
    'add server_parent table'
),
(   'view',
    20191104001,
    'add dictionary schema/definition functions'
),
(   'view',
    20191104002,
    'add resolveRuntimeTo.. functions'
),
(   'filter',
    20191126001,
    'fully re-design filter schema'
),
(   'asset',
    20200908001,
    'add tables for socket in schema asset'
),
(   'ix',
    20200914001,
    'add endpoint tables'
),
(   'filter',
    20200914001,
    'add endpoint as filter-able entity'
),
(   'yp',
    20200915001,
    'rename schema tosm to yp'
),
(   'yp',
    20200915002,
    'add service tables'
),
(   'bulk',
    20200915001,
    'rename instance data table'
),
(   'asset',
    20200915001,
    'remove orchestration environments as possible socket parent'
),
(   'ix',
    20200915001,
    'rename technical system services to technical services'
),
(   'yp',
    '20200917001',
    'cleanup of relationship tables'
),
(   'asset',
    '20200917001',
    'cleanup of relationship tables'
),
(   'ix',
    20200917001,
    'cleanup of relationship tables'
),
(   'asset',
    20201015001,
    'add container entity support'
),
(   'yp',
    20201016001,
    'switch from service linking 1:n to mapping n:m'
),
(   'filter',
    20201016001,
    'add container as filter-able entity'
),
(   'meta',
    20201210001,
    'add meta.attribute registry table'
),
(   'inventory',
    20210105001,
    'add new namespace: inventory'
),
(   'meta',
    20210105001,
    'add inventory information'
),
(   'ix',
    20210726001,
    'add inventory information'
),
(   'yp',
    20210726001,
    'add inventory information'
),
(   'asset',
    20210726001,
    'add inventory information'
),
(   'bulk',
    20210727001,
    'add inventory information'
),
(   'filter',
    20210727001,
    'add inventory information'
),
(   'inventory',
    20210727001,
    'make user attributes optional'
),
(   'inventory',
    20210727002,
    'make inventory names unique per library'
),
(   'inventory',
    20210824001,
    'add user system~nobody'
),
(   'view',
    20220202001,
    'update resolveRuntimeTo.. functions'
),
(   'asset',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'bulk',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'filter',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'inventory',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'ix',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'meta',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'public',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'yp',
    20220204001,
    'update default value for timestamp columns with timezone'
),
(   'view',
    20220211001,
    'update view.resolve..To.. functions'
),
(   'view',
    20220225001,
    'update view.resolve..To..At functions'
),
(   'abstract',
    20220915999,
    'modelupdate'
),
(   'meta',
    20220915999,
    'modelupdate'
),
(   'bulk',
    20220915999,
    'modelupdate'
),
(   'inventory',
    20220915999,
    'modelupdate'
),
(   'yp',
    20220915999,
    'modelupdate'
),
(   'asset',
    20220915999,
    'modelupdate'
),
(   'filter',
    20220915999,
    'modelupdate'
),
(   'view',
    20220915999,
    'modelupdate'
),
(   'ix',
    20220915999,
    'modelupdate'
),
(   'production',
    20220915999,
    'modelupdate'
)
;
