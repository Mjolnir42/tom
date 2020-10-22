BEGIN;
  ALTER TABLE ix.logical_component_subgroup RENAME TO deployment_group;
  ALTER TABLE ix.deployment_group RENAME COLUMN subgroupID TO groupID;
  ALTER INDEX __pk_ixlcs RENAME TO __pk_ixdg;
  ALTER TABLE ix.deployment_group RENAME CONSTRAINT __fk_ixlcs_dictID TO __fk_ixdg_dictID;
  ALTER TABLE ix.deployment_group RENAME CONSTRAINT __ixlcs_fk_origin TO __ixdg_fk_origin;

  ALTER TABLE ix.logical_component_subgroup_attribute RENAME TO deployment_group_attr_values;
  ALTER TABLE ix.deployment_group_attr_values RENAME COLUMN subgroupID TO groupID;
  ALTER TABLE ix.deployment_group_attr_values RENAME CONSTRAINT __fk_ixlcsa_subID TO __fk_ixdgav_grpID;
  ALTER TABLE ix.deployment_group_attr_values RENAME CONSTRAINT __fk_ixlcsa_attrID TO __fk_ixdgav_attrID;
  ALTER TABLE ix.deployment_group_attr_values RENAME CONSTRAINT __fk_ixlcsa_dictID TO __fk_ixdgav_dictID;
  ALTER TABLE ix.deployment_group_attr_values RENAME CONSTRAINT __fk_ixlcsa_uq_dct TO __fk_ixdgav_uq_dct;
  ALTER TABLE ix.deployment_group_attr_values RENAME CONSTRAINT __fk_ixlcsa_uq_att TO __fk_ixdgav_uq_att;
  ALTER TABLE ix.deployment_group_attr_values RENAME CONSTRAINT __ixlcsa_temporal TO __ixdgav_temporal;

  ALTER TABLE ix.logical_component_subgroup_unique RENAME TO deployment_group_attr_uniq_values;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME COLUMN subgroupID TO groupID;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __fk_ixlcsq_subID TO __fk_ixdgqv_grpID;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __fk_ixlcsq_attrID TO __fk_ixdgqv_attrID;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __fk_ixlcsq_dictID TO __fk_ixdgqv_dictID;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __fk_ixlcsq_uq_dct TO __fk_ixdgqv_uq_dct;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __fk_ixlcsq_uq_att TO __fk_ixdgqv_uq_att;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __ixlcsq_temporal TO __ixdgqv_temporal;
  ALTER TABLE ix.deployment_group_attr_uniq_values RENAME CONSTRAINT __ixlcsq_temp_uniq TO __ixdgqv_temp_uniq;

  ALTER TABLE ix.mapping_logical_component_subgroup RENAME TO mapping_deployment_group;
  ALTER TABLE ix.mapping_deployment_group RENAME COLUMN subgroupID TO groupID;
  ALTER TABLE ix.mapping_deployment_group RENAME CONSTRAINT __fk_ixmlcs_grpID TO __fk_ixmdg_grpID;
  ALTER TABLE ix.mapping_deployment_group RENAME CONSTRAINT __fk_ixmlcs_techID TO __fk_ixmdg_techID;
  ALTER TABLE ix.mapping_deployment_group RENAME CONSTRAINT __ixmlcs_temporal TO __ixmdg_temporal;

  ALTER TABLE ix.mapping_functional_component RENAME COLUMN subgroupID TO groupID;

  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'ix', 20191014001, 'rename logical_component_subgroup to deployment_group');
COMMIT;
